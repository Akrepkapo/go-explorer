/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package controllers

import (
	//"encoding/json"
	"fmt"
	"strconv"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Get_FindDatabase(c *gin.Context) {

	ret := &Response{}
	req := &DataBaseFind{}
	rb := DataBaseRespone{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		ret.ReturnFailureString(err.Error())
		JsonResponse(c, ret)
		return
	}

	rb.Cmd = req.Cmd
	rb.Current_page = req.Current_page
	rb.Table_name = req.Table_name
	rb.NodePosition = req.NodePosition
	rb.Page_size = req.Page_size

	switch rb.Cmd {
	case "001":
		rets, err := GetDBinfo()
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		} else {
			rb.Total = len(*rets)
			rb.Data = rets
			ret.Return(rb, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
	case "002":
		total, rets, err := models.GetNodeALLTable(req.NodePosition, req.Current_page, req.Page_size, req.Order, req.Table_name)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		} else {
			rb.Page_size = len(rets)
			rb.Total = int(total)
			rb.Data = rets
			ret.Return(rb, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
	case "003":

		rets, err := models.GetNodeAllColumnTypes(req.NodePosition, req.Table_name)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		} else {
			rb.Total = len(rets)
			rb.Data = rets
			ret.Return(rb, CodeSuccess)
			JsonResponse(c, ret)
			return
		}
	case "004":

		limit := req.Page_size
		order := models.GetNodeTableOrder(req.Order, req.Table_name)
		rb.Order = order
		num, err := models.GetNodeRows(req.NodePosition, req.Table_name)
		if err != nil {
			ret.ReturnFailureString(err.Error())
			JsonResponse(c, ret)
			return
		} else if req.Where == "" {
			list, err := models.GetAll(req.NodePosition, fmt.Sprintf(`select * from "%s" order by %s offset %d`, req.Table_name, order, (req.Current_page-1)*req.Page_size), limit)
			if err != nil {
				ret.ReturnFailureString(err.Error())
				JsonResponse(c, ret)
				return
			} else {
				rb.Data = list
				rb.Page_size = req.Page_size
				rb.Total = int(num)
				ret.Return(rb, CodeSuccess)
				JsonResponse(c, ret)
				return
			}
		} else {
			list, err := models.GetAll(req.NodePosition, fmt.Sprintf(`select * from "%s" where %s order by %s offset %d`, req.Table_name, req.Where, req.Order, (req.Current_page-1)*req.Page_size), limit)
			out := fmt.Sprintf(`select * from %s where %s order by %s offset %d `, req.Table_name, req.Where, order, (req.Current_page-1)*req.Page_size)
			fmt.Println(out)
			if err != nil {
				ret.ReturnFailureString(err.Error())
				JsonResponse(c, ret)
				return
			} else {
				rb.Data = list
				rb.Page_size = req.Page_size
				rb.Total = int(num)
				ret.Return(rb, CodeSuccess)
				JsonResponse(c, ret)
				return
			}
		}

	default:
		ret.ReturnFailureString("cmd err")
		JsonResponse(c, ret)
		return
	}

}

	var (
		ret []DBWebInfo
	)
	dat := models.FullnodesInfo
	dlen := len(dat)
	for i := 0; i < dlen; i++ {
		dat1 := DBWebInfo{
			Id:       strconv.FormatInt(dat[i].NodePosition, 10),
			Nodename: dat[i].Nodename,
			IconUrl:  dat[i].IconUrl,
			Name:     dat[i].Name,
			//Engine:   dat[i].Engine,
			//Version:  dat[i].Version,
		}
		if dat[i].Display {
			ret = append(ret, dat1)
		}
	}
	return &ret, nil
}
