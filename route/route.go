/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package route

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			//c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//OPTIONS
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
func prefix(s string) string {
	return "/api/v2/" + s
}
func Run(host string) (err error) {
	r := gin.Default()
	r.Use(Cors())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "xxxx-api....：V1.0-2020.11.19",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET(prefix("maxblockid"), controllers.Get_maxblockid)
	r.GET(`/api/v2/websocket_token`, controllers.DashboardGetToken)
	r.GET(`/api/v2/honor_nodelist`, controllers.GetHonorNodelists)

	r.GET(`/api/v2/dashboard`, controllers.GetDashboard)
	r.GET(`/api/v2/block_tpslist`, controllers.GetBlockTpsLists)
	r.GET(`/api/v2/blocktransactionlist/:count`, controllers.GetDashboardBlockTransactions)
	r.POST(`/api/v2/common_transaction_search`, controllers.CommonTransactionSearch)

	r.GET(`/api/v2/transaction_detail/:hash`, controllers.Get_TransactionDetails)
	r.GET(`/api/v2/block_detail/:blockid`, controllers.Get_BlockDetails)

	r.POST(`/api/v2/block_detail`, controllers.Get_BlockDetail)

	r.POST(`/api/v2/wallettotal`, controllers.Get_WalletTotal)
	r.POST(`/api/v2/ecosytem_transaction_history`, controllers.Get_EcosytemTranscationhistory)
	r.POST(`/api/v2/database`, controllers.Get_FindDatabase)
	r.POST(`/api/v2/ecosystem_param`, controllers.Get_ecosystem_param)

	r.POST(`/api//v2/get_transaction`, controllers.Get_transaction_block_redis)

	r.GET(`/api/v2/leveldb/:key`, controllers.GetRedisKey)

	r.StaticFS("/api/v2/logo", http.Dir("./logodir"))

	if conf.GetEnvConf().ServerInfo.EnableHttps {
		err = r.RunTLS(host, conf.GetEnvConf().ServerInfo.CertFile, conf.GetEnvConf().ServerInfo.KeyFile)
	} else {
		err = r.Run(host)
	}
	if err != nil {
		log.Errorf("server http/https start failed :%s", err.Error())
		return err
	}

	return nil
}
