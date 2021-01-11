/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"strconv"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/IBAX-io/go-explorer/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type getMaxBlockIDResult struct {
	MaxBlockID int64 `json:"max_block_id"`
}

func Get_maxblockid(c *gin.Context) {
	ret := &Response{}
	bid := getMaxBlockIDResult{}
	if services.Sqlite_MaxBlockid > 0 {
		bid.MaxBlockID = services.Sqlite_MaxBlockid
		ret.Return(bid, CodeSuccess)
		JsonResponse(c, ret)
	} else {
		bk := &models.Block{}
		f, err := bk.GetMaxBlock()
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
		}
		if f {
			bid.MaxBlockID = bk.ID
			ret.Return(bid, CodeSuccess)
			JsonResponse(c, ret)
		}
	}
}

// @tags  block detial
// @Description block detial
// @Summary block detial
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_BlockDetails(c *gin.Context) {
	ret := &Response{}
	blockid := c.Param("blockid")
	bid, err := strconv.ParseInt(blockid, 10, 64)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	var bk models.Block
	fb, err := bk.Get(bid)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	if fb {
		bdt, err := models.GetBlocksDetailedInfoHexByScanOut(&bk)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		}
		var sc models.ScanOut
		sc.BlockDetail = *bdt
		sc.BlockDetail.Header.NodePosition += 1
		sc.BlockDetail.NodePosition += 1
		ret.Return(sc.BlockDetail, CodeSuccess)
		JsonResponse(c, ret)
		return
	}
	ret.ReturnFailureString("not found blockid in db: " + blockid)
	JsonResponse(c, ret)
	return
}

func Get_BlockDetail(c *gin.Context) {
	ret := &Response{}
			bdt, err := models.GetBlocksDetailedInfoHexByScanOut(&bk)
			if err != nil {
				ret.ReturnFailureString(err.Error())
				JsonResponse(c, ret)
				return
			}
			var sc models.ScanOut
			sc.BlockDetail = *bdt
			rs := sc.GetBlockDetialRespones(req.Page, req.Limit)
			rs.Header.NodePosition += 1
			rs.NodePosition += 1
			ret.Return(rs, CodeSuccess)
			JsonResponse(c, ret)
			return
		} else {
			ret.ReturnFailureString("not found blockid in db: " + converter.Int64ToStr(req.Block_id))
			JsonResponse(c, ret)
			return
		}
	}

}
