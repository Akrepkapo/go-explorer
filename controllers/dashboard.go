/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/gin-gonic/gin"
)

func DashboardGetToken(c *gin.Context) {
	ret := &Response{}
	rets, err := services.GetJWTCentToken(1, 60*60)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
	} else {
		//select {
		//case services.SendWebsocketData <- true:
		//default:
		//}
		ret.Return(rets, CodeSuccess)
		JsonResponse(c, ret)
	}
}

func GetDashboard(c *gin.Context) {
	ret := &Response{}
	var scanout models.ScanOut
	rets, err := scanout.GetRedisdashboard()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetBlockTpsLists(c *gin.Context) {
	ret := &Response{}
	rets, err := services.Get_Group_Block_TpsLists()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
	}
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
	return
}

func GetDashboardBlockTransactions(c *gin.Context) {
	ret := &Response{}
	cs := c.Param("count")
	count := converter.StrToInt(cs)
	var scanout models.ScanOut
	rets, err := scanout.GetBlockTransactions(count)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
	} else {
		ret.Return(rets, CodeSuccess)
		JsonResponse(c, ret)
	}

}
