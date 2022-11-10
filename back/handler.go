package main

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/orcastor/orcas/core"
	"github.com/orcastor/orcas/rpc/util"
	"github.com/orcastor/orcas/sdk"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
)

var hanlder = core.NewLocalHandler()

func get(ctx *gin.Context) {
	var req struct {
		BktID int64 `json:"b,omitempty"`
		ID    int64 `json:"i,omitempty"`
	}
	req.BktID, _ = strconv.ParseInt(ctx.Query("b"), 10, 64)
	req.ID, _ = strconv.ParseInt(ctx.Query("i"), 10, 64)
	o, err := hanlder.Get(ctx.Request.Context(), req.BktID, []int64{req.ID})
	if err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	if len(o) <= 0 || o[0].Type != core.OBJ_TYPE_FILE {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	dataID := o[0].DataID
	// 如果不是首版本
	if dataID == 0 {
		os, _, _, err := hanlder.List(ctx.Request.Context(), req.BktID, req.ID, core.ListOptions{
			Type:  core.OBJ_TYPE_VERSION,
			Count: 1,
			Order: "-id",
		})
		if err != nil || len(os) < 1 {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
		dataID = os[0].DataID
	}

	var d *core.DataInfo
	if dataID == core.EmptyDataID {
		d = core.EmptyDataInfo()
	} else {
		d, err = hanlder.GetDataInfo(ctx.Request.Context(), req.BktID, dataID)
		if err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(o[0].Name)))

	var decmpr archiver.Decompressor
	if d.Kind&core.DATA_CMPR_MASK != 0 {
		if d.Kind&core.DATA_CMPR_SNAPPY != 0 {
			decmpr = &archiver.Snappy{}
		} else if d.Kind&core.DATA_CMPR_ZSTD != 0 {
			decmpr = &archiver.Zstd{}
		} else if d.Kind&core.DATA_CMPR_GZIP != 0 {
			decmpr = &archiver.Gz{}
		}
	} else {
		decmpr = &sdk.DummyArchiver{}
	}

	if err = decmpr.Decompress(sdk.NewDataReader(ctx.Request.Context(), hanlder, req.BktID, d, ""), ctx.Writer); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
}
