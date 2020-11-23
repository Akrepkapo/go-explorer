/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	//"encoding/json"

func Get_system_param(c *gin.Context) {

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
	ret, err := model.GetAll(fmt.Sprintf(`select * from "%s" order by %s`, "1_system_parameters", req.Params.Order), 5000)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		rb.Total = int64(len(ret))
		GenResponse(c, req.Head, rb)
	}
}

func Get_ecosystem(c *gin.Context) {
	var (
		ecosystems     []models.Ecosystem
		ecosystemslist []models.EcosystemList
		i              int
	)
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
	err := models.GetALL("1_ecosystems", "", &ecosystems)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		for i = 0; i < len(ecosystems); i++ {
			var ec = models.EcosystemList{
				ID:             ecosystems[i].ID,
				Name:           ecosystems[i].Name,
				IsValued:       ecosystems[i].IsValued,
				Info:           ecosystems[i].Info,
				EmissionAuount: ecosystems[i].EmissionAmount,
				TokenTitle:     ecosystems[i].TokenTitle,
				TypeEmission:   ecosystems[i].TypeEmission,
				TypeWithdraw:   ecosystems[i].TypeWithdraw,
			}

			if ec.ID == 1 {
				ec.TokenTitle = consts.SysEcosytemTitle
			}
			//var (
			//	keys []models.Member
			//)

			id := strconv.FormatInt(ecosystems[i].ID, 10)
			//err = DBConn.Table("1_keys").Where("id = ?", wid).Find(&keys).Error
			count, err := models.GetEcosytem(ecosystems[i].ID)
			if err == nil {
				ec.Member = count
			}

			//list, err := model.GetAll(fmt.Sprintf(`select * from "%s"`, id+"_app_params"), 5000)
			list, err := model.GetAll(fmt.Sprintf(`select * from "%s"`, id+"_parameters"), 5000)
			if err == nil {
				ec.AppParams = list
				//ecosystemslist = append(ecosystemslist, ec)
			}
			ecosystemslist = append(ecosystemslist, ec)
		}

		rb.Data = ecosystemslist
		rb.Total = int64(len(ecosystems))
		GenResponse(c, req.Head, rb)
	}
}

func Get_ecosystem_param(c *gin.Context) {
	var (
		//params []models.AppParam
		rets models.SystemParameterResult
	)
	ret := &Response{}
	req := &EcosytemTranscationHistoryFind{}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	//eid := strconv.FormatInt(req.Ecosystem, 10)
	//err := models.GetALL(eid+"_app_params", "", &rets.Rets)
	//var ap models.AppParam
	//num, rs, err := ap.FindAppParameters(req.Ecosystem, req.Page, req.Limit, req.Search, req.Order)
	//if err != nil {
	//	ret.ReturnFailureString(err.Error())
	//	JsonResponse(c, ret)
	//	return
	//}
	var ap models.SystemParameter
	num, rs, err := ap.FindAppParameters(req.Page, req.Limit, req.Search, req.Order)
	if err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}
	rets.Total = int64(num)
	rets.Page = req.Page
	rets.Limit = req.Limit
	rets.Rets = rs
	ret.Return(rets, CodeSuccess)
	JsonResponse(c, ret)
}

func Get_ecosystem_Keys(c *gin.Context) {
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
	//eid := strconv.FormatInt(req.Params.Current_page, 10)
	key := models.Key{}
	ret, num, err := key.GetKeys(req.Params.Ecosystem, req.Params.CurrentPage, req.Params.PageSize, req.Params.Order)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		rb.Total = num
		GenResponse(c, req.Head, rb)
	}
}

func Get_ecosystem_Key(c *gin.Context) {

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
	key := models.Key{}
	ret, err := key.Get(req.Params.Ecosystem, req.Params.Wallet)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		//rb.Total = len(params)
		GenResponse(c, req.Head, rb)
	}
}

func Get_ecosystemall_Key(c *gin.Context) {

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
	key := models.Key{}
	ret, err := key.GetEcosykey(req.Params.Wallet)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = ret
		//rb.Total = len(params)
		GenResponse(c, req.Head, rb)
	}
}
