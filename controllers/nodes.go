/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

import (
	//"errors"
	//"strconv"
	"encoding/json"

	//"strconv"
	"fmt"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetNodes(c *gin.Context) {

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

	byteNodes += `{"tcp_address":"node5.ibax.one:7081","api_address":"http://node5.ibax.one:7082","city": "china-beijing","icon": "china","key_id":"8226801187119005894","public_key":"ace7333c13567170ed45a300535b3909e6e091350c0dd91b8d8eadd148fc7c59836efb6335e0fab7434d971ebb6539129ae46d539cd92261232dc242261cfcaf"},`
	byteNodes += `{"tcp_address":"hwsh.chain.gs:7081","api_address":"https://hwsh.chain.gs:7082","city": "japan-tokyo","icon": "japan","key_id":"-3232597495798013991","public_key":"0b44d1b8758a32c3d8ecfc3ccc87469dde23d76067d88d81aab4d61952327eb5673d630316b24adac764dfebad1920feb2b6bdfc4cdc44bce6e06dbafc77c084"},`
	byteNodes += `{"tcp_address":"node7.ibax.one:7081","api_address":"http://node7.ibax.one:7082","city": "us-ashburnvirginia","icon": "united_states","key_id":"-4859061282847913966","public_key":"96c3059a1c91f65a4b28d06a307a8839094fa5fc0eae176eca61f75e045fccca2f71fab7b541c4a3587aafd4e8ab0eeeecd7eb65f7561a4be1ca336bf1c83e81"}`
	byteNodes += `]`

	var fs []models.FullNodeCityJSONHex
	err := json.Unmarshal([]byte(byteNodes), &fs)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = fs
		rb.Total = int64(len(fs))
		GenResponse(c, req.Head, rb)
	}

}

func DashboardNodeMap(c *gin.Context) {

	req := &WebRequest{}
	rb := &ResponseBoby{
		Cmd: "001",
		//		Page_size:     "10",
		//		Current_page:    "1",
		//		Total: "100",
		Ret:     "1",
		Retcode: 200,
		Retinfo: "ok",
	}

	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	}

	byteNodes := `[`
	byteNodes += `{"name": "china-guangzhou", "latitude": 23.1367070000, "longitude": 113.2908070000},`
	byteNodes += `{"name": "china-guangzhou", "latitude": 23.2367070000, "longitude": 113.3908070000},`
	byteNodes += `{"name": "china-shanghai", "latitude": 31.2233714132, "longitude": 121.4593505859},`
	byteNodes += `{"name": "china-shanghai", "latitude": 31.4233714132, "longitude": 121.5593505859},`
	byteNodes += `{"name": "china-shanghai", "latitude": 31.5233714132, "longitude": 121.6593505859},`
	byteNodes += `{"name": "china-beijing", "latitude": 39.8097362345, "longitude": 116.6221191406},`
	byteNodes += `{"name": "china-beijing", "latitude": 39.9097362345, "longitude": 116.4221191406},`
	byteNodes += `{"name": "china-chengdu", "latitude": 30.4804397865, "longitude": 104.4899658203},`
	byteNodes += `{"name": "china-chengdu", "latitude": 30.6804397865, "longitude": 104.0899658203},`
	byteNodes += `{"name": "china-hongkong", "latitude": 22.3011673701, "longitude": 114.1815948486},`
	byteNodes += `{"name": "japan-tokyo", "latitude": 35.8467608768, "longitude": 139.6994018555},`
	byteNodes += `{"name": "us-sanfrancisco", "latitude": 37.7749290000, "longitude": -122.4194160000},`
	byteNodes += `{"name": "singapore-singapore", "latitude": 1.8838157763, "longitude": 103.4527587891}`
	byteNodes += `]`

	fmt.Println(byteNodes)
	var fs []NodeMapInfo
	err := json.Unmarshal([]byte(byteNodes), &fs)
	if err != nil {
		rb.Retinfo = err.Error()
		rb.Retcode = 404
		GenResponse(c, req.Head, rb)
	} else {
		rb.Data = fs
		rb.Total = int64(len(fs))
		GenResponse(c, req.Head, rb)
	}

}

func GetHonorNodelists(c *gin.Context) {
	ret := &Response{}
	ret.Return(models.IsDisplay(models.FullnodesInfo), CodeSuccess)
	JsonResponse(c, ret)

}
