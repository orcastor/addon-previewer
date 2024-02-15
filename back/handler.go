package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
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
	"dwg": true,
	"dxf": true,
}

var iworkTypes = map[string]bool{
	"pages":   true,
	"numbers": true,
	"key":     true,
}

var docConvTypes = map[string]string{
	"dwg": "png",
	"dxf": "png",

	"pages":   "html",
	"numbers": "html",
	"key":     "html",

	"dsp":  "pdf",
	"ppt":  "pdf",
	"pptx": "pdf",

	"et":  "xlsx",
	"csv": "xlsx",
	"xls": "xlsx",

	"wps":    "docx",
	"doc":    "docx",
	"txt":    "docx",
	"json":   "docx",
	"toml":   "docx",
	"yaml":   "docx",
	"xml":    "docx",
	"config": "docx",
}

func get(ctx *gin.Context) {
	bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
	id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)
	langID, _ := strconv.ParseInt(ctx.Query("l"), 10, 64) // 预留的语言id

	from := strings.ToLower(ctx.Query("t")) // from无法注入，不在白名单会直接返回
	to := docConvTypes[from]
	if to == "" {
		// 不需要转换格式，那就直接写到http
		if err := writeTo(ctx, bktID, id, ctx.Writer, true); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
		}
		return
	}

	fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, from)) // id无法注入，强制转成数字
	toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, to))

	// 先看转换后文件生成了没有
	if st, err := os.Stat(toPath); err == nil && st.Size() > 0 {
		goto READ_OUTPUT_FILE
	}

	// 下载到临时文件
	if err := download(ctx, bktID, id, fromPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 转换格式
	if cadTypes[from] {
		if err := cad2xConv(fromPath, toPath, langID); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	} else if iworkTypes[from] {
		if err := iwork2htmlConv(fromPath, toPath); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	} else {
		if err := x2tConv(fromPath, toPath); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	}

	// 删除临时文件
	os.Remove(fromPath)

READ_OUTPUT_FILE:
	ft, err := os.Open(toPath)
	if err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
	defer ft.Close()

	if _, err = io.Copy(ctx.Writer, bufio.NewReader(ft)); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
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

func cad2xConv(fromPath, toPath string, langID int64) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "),
		"/opt/cad2x/cad2pdf",
		fromPath,
		"-o", toPath,
		"-ac",
		"-e", langID2CodePage(langID),
		"-f", "simsun",
		"-l", "/opt/x2t/core-fonts")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("cad2xConv error: %+v", out.String())
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
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "),
		"/opt/x2t/x2t",
		fromPath,
		toPath)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("x2tConv error: %+v", out.String())
	}
	return err
}

var icoTypes = map[string]bool{
	"dmg":     true,
	"exe":     true,
	"apk":     true,
	"inf":     true,
	"ini":     true,
	"desktop": true,
	"app":     true,
}

var thumbSupport = map[string]bool{
	"csv":  true,
	"bmp":  true,
	"raw":  true,
	"jpg":  true,
	"jpeg": true,
	"jpe":  true,
	"jfif": true,
	"png":  true,
	"gif":  true,
	"tif":  true,
	"tiff": true,
	"webp": true,

	"mat":  true,
	"pbm":  true,
	"pgm":  true,
	"ppm":  true,
	"pfm":  true,
	"pnm":  true,
	"fits": true,
	"fit":  true,
	"fts":  true,
	"exr":  true,
	"hdr":  true,
	"v":    true,
	"vips": true,
}

var outTypes = map[string]bool{
	"jpg":  true,
	"png":  true,
	"gif":  true,
	"webp": true,
	"ico":  true,
}

func thumb(ctx *gin.Context) {
	bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
	id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)

	w, _ := strconv.ParseInt(ctx.Query("w"), 10, 64)
	h, _ := strconv.ParseInt(ctx.Query("h"), 10, 64)

	from := strings.ToLower(ctx.Query("t")) // from无法注入，不在白名单会直接返回
	if !thumbSupport[from] && !icoTypes[from] {
		util.AbortResponse(ctx, 400, "not supported format")
		return
	}

	// TODO：如果是文档格式，先用x2t获取文档缩略图

	// TODO：如果涉及隐私文件，返回不支持获取缩略图

	to := strings.ToLower(ctx.Query("nt")) // to无法注入，不在白名单会直接返回
	if !outTypes[to] {
		util.AbortResponse(ctx, 400, "not supported format")
		return
	}

	fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, from)) // id无法注入，强制转成数字
	toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%dx%d.%s", id, w, h, to))

	// 先看转换后文件生成了没有
	if st, err := os.Stat(toPath); err == nil && st.Size() > 0 {
		goto READ_OUTPUT_FILE
	}

	// 下载到临时文件
	if err := download(ctx, bktID, id, fromPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 如果是获取图标，用f2ico处理
	if icoTypes[from] {
		f, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err = fico.F2ICO(f, fromPath, f2ico.Config{Width: int(w), Height: int(h)}); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	} else if err := vipsConv(fromPath, toPath, w, h); err != nil {
		// 转换缩略图
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 删除临时文件
	os.Remove(fromPath)

READ_OUTPUT_FILE:
	ft, err := os.Open(toPath)
	if err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
	defer ft.Close()

	if _, err = io.Copy(ctx.Writer, bufio.NewReader(ft)); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
}

func vipsConv(fromPath, toPath string, w, h int64) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "), "/opt/vips/vipsthumbnail", fromPath,
		"--size", fmt.Sprintf("%dx%d", w, h),
		"--smartcrop", "attention",
		"-o", toPath+"[keep=none]")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("vipsConv error: %+v", out.String())
	}
	return err
}

func download(ctx *gin.Context, bktID, id int64, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	defer b.Flush()

	return writeTo(ctx, bktID, id, b, false)
}

func writeTo(ctx *gin.Context, bktID, id int64, writer io.Writer, direct bool) error {
	o, err := hanlder.Get(ctx.Request.Context(), bktID, []int64{id})
	if err != nil {
		return err
	}

	if len(o) <= 0 || o[0].Type != core.OBJ_TYPE_FILE {
		return errors.New("not file")
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
			return err
		}
		if len(os) < 1 {
			return errors.New("no version")
		}
		dataID = os[0].DataID
	}

	var d *core.DataInfo
	if dataID == core.EmptyDataID {
		d = core.EmptyDataInfo()
	} else {
		d, err = hanlder.GetDataInfo(ctx.Request.Context(), bktID, dataID)
		if err != nil {
			return err
		}
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(o[0].Name)))

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
