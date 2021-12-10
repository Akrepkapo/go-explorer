/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package controllers

import (
	"net/http"

	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func JsonResponse(c *gin.Context, body *Response) {
	c.JSON(http.StatusOK, body)
}

//GenResponse genrate reponse ,json format
func GenResponse(c *gin.Context, head *RequestHead, body *ResponseBoby) {
	c.JSON(http.StatusOK, gin.H{
		"body": body,
		"head": head,
	})
}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Common_search(c *gin.Context) {

	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd:     "001",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	rb.PageSize = req.Params.PageSize
	rb.CurrentPage = req.Params.CurrentPage
	if len(req.Params.Hash) == 64 {
		//
		ret, err := services.Get_Group_Block_Detail_hash(req.Params.Hash)
		if err == nil && ret.Header.BlockID > 0 {
			rb.Data = ret
			rb.RetDataType = "1"
			GenResponse(c, req.Head, rb)

		} else {
			logourl := conf.GetEnvConf().Url.URL
			ret, err1 := services.Get_transaction_Hash(logourl, req.Params.Hash)
			if err1 != nil {
				rb.Retinfo = err1.Error()
				rb.Retcode = 404
				GenResponse(c, req.Head, rb)
			} else {
				rb.Data = ret
				rb.RetDataType = "2"
				rb.Total = int64(len(*ret))
				GenResponse(c, req.Head, rb)
			}

		}
	} else {
		//
		ret, _, err1 := services.GetBlockDetailed(req.Params.Block_id)
		if err1 != nil {
			rb.Retinfo = err1.Error()
			rb.Retcode = 404
			GenResponse(c, req.Head, rb)
		} else {
			rb.Data = ret
			rb.RetDataType = "1"
			GenResponse(c, req.Head, rb)
		}
	}
}

// @tags  ecosystem
// @Description ecosystem
// @Summary ecosystem
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
//TODO:20210901 This interface needs to be changed to read data from redis to fetch from the database
func GetRedisKey(c *gin.Context) {
	ret := &Response{}
	cs := c.Param("key")
	count := converter.StrToInt64(cs)
	var scanout models.ScanOut
	f, err := scanout.Get_Redis(count)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		JsonResponse(c, ret)
		return
	}
	if f {
		ret.Return(scanout, CodeSuccess)
		JsonResponse(c, ret)
		return
	} else {
		ret.ReturnFailureString("not found key in redis: " + cs)
		JsonResponse(c, ret)
		return
	}

}
