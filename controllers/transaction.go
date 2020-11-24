/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_transaction(c *gin.Context) {
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
	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := services.Get_Group_TransactionStatus(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err == nil && ret != nil {
		rb.Data = ret
		rb.Total = num
		//rb.Page_size = req.Params.Page_size
		//rb.Current_page = req.Params.Current_page
		GenResponse(c, req.Head, rb)
	} else {
		if err != nil {
			rb.Retinfo = err.Error()
		}
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}
}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_transaction_history(c *gin.Context) {
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

	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := services.Get_Group_TransactionHistory(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err == nil {
		rb.Data = ret
		rb.Total = num

		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_transaction_block(c *gin.Context) {
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

	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	ret, num, err := models.Get_Group_TransactionBlock(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err == nil && ret != nil {
		rb.Data = ret
		rb.Total = num
		GenResponse(c, req.Head, rb)
	} else {
		if err != nil {
			rb.Retinfo = err.Error()
		}
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_transaction_details(c *gin.Context) {

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
		//rb.Body.Retinfo = err.Error()
		GenResponse(c, req.Head, rb)
	}

	logourl := conf.GetEnvConf().Url.URL
	ret, err := services.Get_transaction_Hash(logourl, req.Params.Hash)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		rb.PageSize = req.Params.PageSize
		rb.CurrentPage = req.Params.CurrentPage
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_TransactionDetails(c *gin.Context) {

	ret := &Response{}
	hash := c.Param("hash")
	logourl := conf.GetEnvConf().Url.URL
	rets, err := services.Get_transaction_Hash(logourl, hash)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_Find_history(c *gin.Context) {
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
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_Find_Ecosytemhistory(c *gin.Context) {
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

	rb.Ecosystem = req.Params.Ecosystem
	rb.PageSize = req.Params.PageSize
	rb.CurrentPage = req.Params.CurrentPage
	rb.RetDataType = req.Params.SearchType

	ret, num, total, err := services.Get_Group_TransactionEcosytemWallet(req.Params.Ecosystem, req.Params.CurrentPage, req.Params.PageSize, req.Params.Wallet, req.Params.SearchType)
	if err == nil {
		rb.Data = ret
		rb.Total = num
		rb.Sum = total
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_EcosytemTranscationhistory(c *gin.Context) {
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ts := &models.History{}
	rets, err := ts.GetEcosytemTransactionWallets(req.Ecosystem, req.Page, req.Limit, req.Wallet, req.Search)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}

// @tags
// @Description
// @Summary
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_Find_Wallethistory(c *gin.Context) {
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

	//rb.Page_size = req.Params.Page_size
	rb.Ecosystem = req.Params.Ecosystem
	rb.Wallet = req.Params.Wallet

	ret, err := services.Get_Group_WalletHistory(req.Params.Ecosystem, req.Params.Wallet)
	if err == nil {
		rb.Data = ret
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary Find a list of all currencies under the account
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /api/get_wallettotal [post]
func Get_Wallet_Total(c *gin.Context) {
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

	//rb.Page_size = req.Params.Page_size
	//rb.Ecosystem = req.Params.Ecosystem
	rb.Wallet = req.Params.Wallet
	rb.CurrentPage = req.Params.CurrentPage
	rb.PageSize = req.Params.PageSize
	rb.Order = req.Params.Order

	total, page, ret, err := services.Get_Group_Wallet_Total(req.Params.CurrentPage, req.Params.PageSize, req.Params.Order, req.Params.Wallet)
	if err == nil {
		rb.Data = ret
		rb.Total = total
		rb.PageSize = page
		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}

// @tags
// @Description
// @Summary Find a list of all currencies under the account
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /api/get_wallettotal [post]
func Get_WalletTotal(c *gin.Context) {

	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	logourl := conf.GetEnvConf().Url.URL
	ts := &models.Key{}
	rets, err := ts.GetWalletTotal(req.Page, req.Limit, req.Order, req.Wallet, logourl)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

// @tags  common_transaction_search
// @Description common_transaction_search
// @Summary common_transaction_search
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func CommonTransactionSearch(c *gin.Context) {

	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	ts := &models.BlockTxDetailedInfoHex{}
	rets, err := ts.GetCommonTransactionSearch(req.Page, req.Limit, req.Search, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)

}
func Get_transaction_block_redis(c *gin.Context) {
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	rets, _, err := models.GetTransactionBlockFromRedis()
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

// @tags  transaction history
// @Description transaction history
// @Summary transaction history
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true},"message":"success"}}"
// @Router /auth/admin/{id} [get]
func Get_Transaction_queue(c *gin.Context) {
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
	rb.RetDataType = req.Params.SearchType

	ret, num, err := models.GetTransactionpages(req.Params.CurrentPage, req.Params.PageSize)
	if err == nil {
		rb.Data = ret
		rb.Total = num

		GenResponse(c, req.Head, rb)
	} else {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

}
