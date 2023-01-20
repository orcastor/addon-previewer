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

var extConvertMap = map[string]string{
	"ppt":  "pdf",
	"pptx": "pdf",

	"csv": "xlsx",
	"xls": "xlsx",

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

	from := ctx.Query("t") // from无法注入，不在白名单会直接返回
	to := extConvertMap[from]

	// 不需要转换格式，那就直接写到http
	if to == "" {
		if err := writeTo(ctx, bktID, id, ctx.Writer, true); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
		}
		return
	}

	fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, from)) // id无法注入，强制转成数字
	toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, to))

	// 先看转换后文件生成了没有
	if st, err := os.Stat(toPath); err == nil && st.Size() > 0 {
		goto READ_TO_FILE
	}

	// 下载到临时文件
	if err := download(ctx, bktID, id, fromPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 转换格式
	if err := convert(fromPath, toPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 删除临时文件
	os.Remove(fromPath)

READ_TO_FILE:
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

func convert(fromPath, toPath string) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "), "/opt/x2t/x2t", fromPath, toPath)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("convert error: %+v", out.String())
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
