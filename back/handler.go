package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/orcastor/fico"
	"github.com/orcastor/iwork-converter/iwork2html"
	"github.com/orcastor/orcas/core"
	"github.com/orcastor/orcas/rpc/util"
	"github.com/orcastor/orcas/sdk"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
	"github.com/mholt/archiver/v3"
	arv4 "github.com/orcastor/archiver/v4"
)

var ORCAS_DOCKER_EXEC = os.Getenv("ORCAS_DOCKER_EXEC")
var ORCAS_CACHE = os.Getenv("ORCAS_CACHE")

func init() {
	if ORCAS_CACHE == "" {
		ORCAS_CACHE = "/tmp/orcas_cache"
	}
	os.MkdirAll(ORCAS_CACHE, 0766)
}

var hanlder = core.NewLocalHandler()

var cadTypes = map[string]bool{
	".dwg": true,
	".dxf": true,
}

var iworkTypes = map[string]bool{
	".pages":   true,
	".numbers": true,
	".key":     true,
}

var docConvTypes = map[string]string{
	".dwg": "png",
	".dxf": "png",

	".pages":   "html",
	".numbers": "html",
	".key":     "html",

	".dsp":  "pdf",
	".ppt":  "pdf",
	".pptx": "pdf",

	".et":  "xlsx",
	".csv": "xlsx",
	".xls": "xlsx",

	".wps":    "docx",
	".doc":    "docx",
	".txt":    "docx",
	".json":   "docx",
	".toml":   "docx",
	".yaml":   "docx",
	".xml":    "docx",
	".config": "docx",
	".xps":    "docx",
}

func hash(r string) uint64 {
	h := md5.New()
	h.Write([]byte(r))
	return binary.BigEndian.Uint64(h.Sum(nil)[4:12])
}

func get(ctx *gin.Context) {
	if err := func() error {
		bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
		id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)
		r := ctx.Query("r")
		langID, _ := strconv.ParseInt(ctx.Query("l"), 10, 64) // 预留的语言id

		d, name, err := getData(ctx, bktID, id)
		if err != nil {
			return err
		}

		oldfrom := strings.ToLower(filepath.Ext(name))
		if r != "" {
			name = filepath.Base(r)
		}

		from := strings.ToLower(filepath.Ext(name)) // from无法注入，不在白名单会直接返回
		to := docConvTypes[from]
		if to == "" {
			// 不需要转换格式，那就直接写到http
			ctx.Header("Content-Type", "application/octet-stream")
			ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(name)))
			return writeTo(ctx, bktID, d, ctx.Writer, true)
		}

		fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d%s", id, oldfrom)) // id无法注入，强制转成数字
		toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%0x.%s", id, hash(r), to))

		// 先看转换后文件生成了没有
		if fileExists(toPath) {
			goto READ_OUTPUT_FILE
		}

		// 下载到临时文件
		if err := download(ctx, bktID, d, fromPath); err != nil {
			return err
		}

		if r != "" {
			// 再下载文件
			if err = zipHandle(ctx, fromPath, r, langID, func(f fs.File) error {
				fromPath = filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%0x%s", id, hash(r), from))
				w, err := create(fromPath)
				if err != nil {
					return err
				}
				defer w.Close()

				_, err = io.Copy(w, f)
				return err
			}); err != nil {
				return err
			}
		}

		// 转换格式
		if cadTypes[from] {
			if err := cad2xConv(fromPath, toPath, langID); err != nil {
				return err
			}
		} else if iworkTypes[from] {
			if err := iwork2htmlConv(fromPath, toPath); err != nil {
				return err
			}
		} else {
			if err := x2tConv(fromPath, toPath); err != nil {
				return err
			}
		}

		// 删除临时文件
		os.Remove(fromPath)

	READ_OUTPUT_FILE:
		ft, err := os.Open(toPath)
		if err != nil {
			return err
		}
		defer ft.Close()

		st, err := ft.Stat()
		if err != nil {
			return err
		}

		http.ServeContent(ctx.Writer, ctx.Request, filepath.Base(toPath), st.ModTime(), ft)
		return nil
	}(); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
	}
}

func langID2CodePage(langID int64) string {
	switch langID {
	case 1054:
		return "ANSI_874" // 泰语
	case 1041:
		return "ANSI_932" // 日语
	case 0, 2052:
		return "ANSI_936" // 简体中文
	case 1042:
		return "ANSI_949" // 韩语
	case 1028:
		return "ANSI_950" // 繁体中文
	case 1029:
		return "ANSI_1250" // 中欧: 保加利亚语
	case 1049, 1062, 1087, 1104, 1121, 2072, 2073:
		return "ANSI_1251" // 西里尔语: 俄语, 白俄罗斯语, 马其顿语, 塞尔维亚语, 斯洛文尼亚语, 乌兹别克语, 波斯尼亚语
	case 1032:
		return "ANSI_1253" // 希腊语
	case 1055:
		return "ANSI_1254" // 土耳其语
	case 1037:
		return "ANSI_1255" // 希伯来语
	case 1025:
		return "ANSI_1256" // 阿拉伯语
	case 1066:
		return "ANSI_1257" // 波罗的海语
	case 1068:
		return "ANSI_1258" // 越南语
	}
	return "ANSI_1252" // 西欧: 英语
}

func codePage2Charset(codePage string) string {
	switch codePage {
	case "ANSI_874": // Thai
		return "windows874"
	case "ANSI_932": // Japanese
		return "shiftjis"
	case "ANSI_936": // UnifiedChinese
		return "gb18030"
	case "ANSI_949": // Korean
		return "euckr"
	case "ANSI_950": // TradChinese
		return "big5"
	case "ANSI_1250": // CentralEurope
		return "iso8859_2"
	case "ANSI_1251": // Cyrillic  or ISO 8859_5
		return "koi8u"
	case "ANSI_1253": // Greek
		return "iso8859_7"
	case "ANSI_1254": // Turkish
		return "iso8859_9"
	case "ANSI_1255": // Hebrew
		return "iso8859_8"
	case "ANSI_1256": // Arabic
		return "iso8859_6"
	case "ANSI_1257": // Baltic
		return "iso8859_13"
	case "ANSI_1258": // Vietnam
		return "iso8859_1"
	}
	return "iso8859_1" // WesternEurope
}

func cad2xConv(fromPath, toPath string, langID int64) error {
	cmds := []string{
		"/opt/cad2x/cad2x",
		fromPath,
		"-o", toPath,
		"-ac",
		"-e", langID2CodePage(langID),
		"-f", "simsun",
		"-l", "/opt/x2t/core-fonts"}
	if ORCAS_DOCKER_EXEC != "" {
		cmds = append(strings.Split(ORCAS_DOCKER_EXEC, " "), cmds...)
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("cad2xConv error: %s %+v", strings.Join(cmds, " "), out.String())
	}
	return err
}

func iwork2htmlConv(fromPath, toPath string) error {
	if err := iwork2html.Convert(fromPath, toPath); err != nil {
		elog.Errorf("iwork2htmlConv error: %+v", err)
		return err
	}
	return nil
}

func x2tConv(fromPath, toPath string) error {
	cmds := []string{
		"/opt/x2t/x2t",
		fromPath,
		toPath}
	if ORCAS_DOCKER_EXEC != "" {
		cmds = append(strings.Split(ORCAS_DOCKER_EXEC, " "), cmds...)
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	elog.Errorf("x2tConv error: %s %+v", strings.Join(cmds, " "), out.String())
	return err
}

var icoTypes = map[string]bool{
	".ico":     true,
	".icns":    true,
	".exe":     true,
	".dll":     true,
	".mui":     true,
	".mun":     true,
	".apk":     true,
	".ipa":     true,
	".inf":     true,
	".ini":     true,
	".desktop": true,
	".app":     true,
}

var thumbSupport = map[string]bool{
	".csv":  true,
	".bmp":  true,
	".raw":  true,
	".jpg":  true,
	".jpeg": true,
	".jpe":  true,
	".jfif": true,
	".png":  true,
	".gif":  true,
	".tif":  true,
	".tiff": true,
	".webp": true,

	".mat":  true,
	".pbm":  true,
	".pgm":  true,
	".ppm":  true,
	".pfm":  true,
	".pnm":  true,
	".fits": true,
	".fit":  true,
	".fts":  true,
	".exr":  true,
	".hdr":  true,
	".v":    true,
	".vips": true,
}

var showTypes = map[string]bool{
	".docx": true,
	".pptx": true,
	".pdf":  true,
	".xlsx": true,
}

var outMimeTypes = map[string]string{
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"ico":  "image/x-icon",
}

func thumb(ctx *gin.Context) {
	if err := func() error {
		bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
		id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)
		r := ctx.Query("r")
		langID, _ := strconv.ParseInt(ctx.Query("l"), 10, 64) // 预留的语言id

		w, _ := strconv.ParseInt(ctx.Query("w"), 10, 64)
		h, _ := strconv.ParseInt(ctx.Query("h"), 10, 64)

		d, name, err := getData(ctx, bktID, id)
		if err != nil {
			return err
		}

		oldfrom := strings.ToLower(filepath.Ext(name))
		if r != "" {
			name = filepath.Base(r)
		}

		from := strings.ToLower(filepath.Ext(name)) // from无法注入，不在白名单会直接返回
		if !thumbSupport[from] && !icoTypes[from] && docConvTypes[from] == "" && !showTypes[from] {
			util.AbortResponse(ctx, 400, "not supported format")
			return nil
		}

		// TODO：如果涉及隐私文件，返回不支持获取缩略图

		to := strings.ToLower(ctx.Param("to")) // to无法注入，不在白名单会直接返回
		mimeType := outMimeTypes[to]
		if mimeType == "" {
			util.AbortResponse(ctx, 400, "not supported format")
			return nil
		}

		fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d%s", id, oldfrom)) // id无法注入，强制转成数字
		toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%0x_%dx%d.%s", id, hash(r), w, h, to))

		ctx.Header("Content-Type", mimeType)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filepath.Base(toPath)))

		// 先看转换后文件生成了没有
		if fileExists(toPath) {
			goto READ_OUTPUT_FILE
		}

		// 下载到临时文件
		if err := download(ctx, bktID, d, fromPath); err != nil {
			return err
		}

		if r != "" {
			// 再下载文件
			if err = zipHandle(ctx, fromPath, r, langID, func(f fs.File) error {
				fromPath = filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%0x%s", id, hash(r), from))
				w, err := create(fromPath)
				if err != nil {
					return err
				}
				defer w.Close()

				_, err = io.Copy(w, f)
				return err
			}); err != nil {
				return err
			}
		}

		// 如果是文档格式，先用x2t获取文档缩略图
		if iworkTypes[from] {
			r, err := zip.OpenReader(fromPath)
			if err != nil {
				return err
			}
			defer r.Close()

			var previewJPG *zip.File
			for _, f := range r.File {
				switch {
				case strings.Contains(f.Name, "preview.jpg"):
					previewJPG = f
				}
			}

			if previewJPG == nil {
				return errors.New("no preview.jpg found")
			}

			toPathJPG := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.jpg", id))
			f, err := create(toPathJPG)
			if err != nil {
				return err
			}
			defer f.Close()

			rc, err := previewJPG.Open()
			if err != nil {
				return err
			}

			defer rc.Close()
			if _, err = io.Copy(f, rc); err != nil {
				return err
			}

			fromPath = toPathJPG
		} else if _, ok := docConvTypes[from]; ok || showTypes[from] {
			toPathPNG := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.png", id))

			if cadTypes[from] {
				if err := cad2xConv(fromPath, toPathPNG, langID); err != nil {
					return err
				}

				fromPath = toPathPNG
			} else if err := x2tConv(fromPath, toPathPNG); err != nil {
				return err
			}

			fromPath = toPathPNG
		}

		// 如果是获取图标，用f2ico处理
		if icoTypes[from] {
			f, err := create(toPath)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if err = fico.F2ICO(f, fromPath, fico.Config{Format: to, Width: int(w), Height: int(h)}); err != nil {
				return err
			}
		} else if err := vipsConv(fromPath, toPath, w, h); err != nil {
			// 转换缩略图
			return err
		}

		// 删除临时文件
		os.Remove(fromPath)

	READ_OUTPUT_FILE:
		ft, err := os.Open(toPath)
		if err != nil {
			return err
		}
		defer ft.Close()

		st, err := ft.Stat()
		if err != nil {
			return err
		}

		http.ServeContent(ctx.Writer, ctx.Request, filepath.Base(toPath), st.ModTime(), ft)
		return nil
	}(); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
	}
}

func vipsConv(fromPath, toPath string, w, h int64) error {
	cmds := []string{"/opt/vips/vipsthumbnail", fromPath,
		"--size", fmt.Sprintf("%dx%d", w, h),
		"--smartcrop", "attention",
		"-o", toPath + "[keep=none]"}
	if ORCAS_DOCKER_EXEC != "" {
		cmds = append(strings.Split(ORCAS_DOCKER_EXEC, " "), cmds...)
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("vipsConv error: %s %+v", strings.Join(cmds, " "), out.String())
	}
	return err
}

func create(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

func fileExists(filePath string) bool {
	st, err := os.Stat(filePath)
	return err == nil && st.Size() > 0
}

func download(ctx *gin.Context, bktID int64, d *core.DataInfo, path string) error {
	f, err := create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeTo(ctx, bktID, d, f, false)
}

func getData(ctx *gin.Context, bktID, id int64) (*core.DataInfo, string, error) {
	o, err := hanlder.Get(ctx.Request.Context(), bktID, []int64{id})
	if err != nil {
		return nil, "", err
	}

	if len(o) <= 0 || o[0].Type != core.OBJ_TYPE_FILE {
		return nil, "", errors.New("not file")
	}

	dataID := o[0].DataID
	// 如果不是首版本
	if dataID == 0 {
		os, _, _, err := hanlder.List(ctx.Request.Context(), bktID, id, core.ListOptions{
			Type:  core.OBJ_TYPE_VERSION,
			Count: 1,
			Order: "-id",
		})
		if err != nil {
			return nil, "", err
		}
		if len(os) < 1 {
			return nil, "", errors.New("no version")
		}
		dataID = os[0].DataID
	}

	var d *core.DataInfo
	if dataID == core.EmptyDataID {
		d = core.EmptyDataInfo()
	} else {
		d, err = hanlder.GetDataInfo(ctx.Request.Context(), bktID, dataID)
		if err != nil {
			return nil, "", err
		}
	}
	return d, o[0].Name, err
}

func writeTo(ctx *gin.Context, bktID int64, d *core.DataInfo, writer io.Writer, direct bool) error {
	acceptEncoding := ctx.Request.Header["Accept-Encoding"]

	var decmpr archiver.Decompressor
	if d.Kind&core.DATA_CMPR_MASK != 0 {
		if d.Kind&core.DATA_CMPR_SNAPPY != 0 {
			decmpr = &archiver.Snappy{}
		} else if d.Kind&core.DATA_CMPR_ZSTD != 0 {
			decmpr = &archiver.Zstd{}
		} else if d.Kind&core.DATA_CMPR_GZIP != 0 {
			if direct && len(acceptEncoding) > 0 && strings.Contains(acceptEncoding[0], "gzip") {
				// 如果浏览器支持gzip，直接返回原始数据
				decmpr = &sdk.DummyArchiver{}
				ctx.Header("Content-Encoding", "gzip")
			} else {
				decmpr = &archiver.Gz{}
			}
		} else if d.Kind&core.DATA_CMPR_BR != 0 {
			if direct && len(acceptEncoding) > 0 && strings.Contains(acceptEncoding[0], "br") {
				// 如果浏览器支持br，直接返回原始数据
				decmpr = &sdk.DummyArchiver{}
				ctx.Header("Content-Encoding", "br")
			} else {
				decmpr = &archiver.Brotli{}
			}
		}
	} else {
		decmpr = &sdk.DummyArchiver{}
	}

	return decmpr.Decompress(sdk.NewDataReader(ctx.Request.Context(), hanlder, bktID, d, ""), writer)
}

type ObjectInfo struct {
	MTime int64  `json:"m,omitempty"` // 更新时间，秒级时间戳
	Type  int    `json:"t,omitempty"` // 对象类型，-1: malformed, 0: none, 1: dir, 2: file, 3: version, 4: preview(thumb/m3u8/pdf)
	Name  string `json:"n,omitempty"` // 对象名称
	Size  int64  `json:"s,omitempty"` // 对象的大小，目录的大小是子对象数，文件的大小是最新版本的字节数
}

var archiveTypes = map[string]bool{
	".zip": true,
	".rar": true,
	".7z":  true,
	".sz":  true,
	".xz":  true,
	".gz":  true,
	".tar": true,
	".bz2": true,
	".br":  true,
	".zz":  true,
	".zst": true,
	".lz4": true,
}

func list(ctx *gin.Context) {
	if err := func() error {
		var req struct {
			BktID  int64  `json:"b,omitempty"`
			PID    int64  `json:"p,omitempty"`
			Route  string `json:"r,omitempty"`
			LangID int64  `json:"l,omitempty"`
			core.ListOptions
		}
		ctx.BindJSON(&req)

		d, name, err := getData(ctx, req.BktID, req.PID)
		if err != nil {
			return err
		}

		from := strings.ToLower(filepath.Ext(name)) // from无法注入，不在白名单会直接返回
		if !archiveTypes[from] {
			util.AbortResponse(ctx, 400, "not supported format")
			return nil
		}

		toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d%s", req.PID, from)) // id无法注入，强制转成数字

		// 先看转换后文件生成了没有
		if !fileExists(toPath) {
			// 下载到临时文件
			if err := download(ctx, req.BktID, d, toPath); err != nil {
				return err
			}
		}

		delimiter := ""
		var o []ObjectInfo
		if err = zipHandle(ctx, toPath, req.Route, req.LangID, func(f fs.File) error {
			if dir, ok := f.(fs.ReadDirFile); ok {
				// 0 gets all entries, but you can pass > 0 to paginate
				entries, err := dir.ReadDir(0)
				if err != nil {
					return err
				}
				for _, e := range entries {
					st, err := e.Info()
					if err != nil {
						continue
					}
					if st.IsDir() {
						o = append(o, ObjectInfo{
							MTime: st.ModTime().Unix(),
							Type:  1,
							Name:  st.Name(),
						})
					} else {
						o = append(o, ObjectInfo{
							MTime: st.ModTime().Unix(),
							Type:  2,
							Name:  st.Name(),
							Size:  st.Size(),
						})
					}
					delimiter = st.Name()
				}
			}
			return nil
		}); err != nil {
			return err
		}

		ctx.AbortWithStatusJSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"c": len(o),
				"d": delimiter,
				"o": o,
			},
		})
		return nil
	}(); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
	}
}

func zipHandle(ctx *gin.Context, path, route string, langID int64, h func(f fs.File) error) error {
	ft, err := os.Open(path)
	if err != nil {
		return err
	}
	defer ft.Close()

	st, err := ft.Stat()
	if err != nil {
		return err
	}

	format, _, err := arv4.Identify(filepath.Base(path), ft)
	if err != nil && !errors.Is(err, arv4.ErrNoMatch) {
		return err
	}

	var fsys arv4.ArchiveFS
	if format != nil {
		switch ff := format.(type) {
		case arv4.Archival:
			switch fff := ff.(type) {
			case arv4.Zip:
				fff.TextEncoding = codePage2Charset(langID2CodePage(langID))
				fff.ContinueOnError = true
				ff = fff
			case arv4.Rar:
				fff.TextEncoding = codePage2Charset(langID2CodePage(langID))
				fff.ContinueOnError = true
				ff = fff
			}
			fsys = arv4.ArchiveFS{Stream: io.NewSectionReader(ft, 0, st.Size()), Format: ff, Context: ctx}
		default:
			util.AbortResponse(ctx, 400, "not supported format")
			return nil
		}
	}

	if route == "" {
		route = "."
	}
	route = strings.TrimPrefix(route, "./")

	f, err := fsys.Open(route)
	if err != nil {
		return err
	}
	defer f.Close()

	return h(f)
}
