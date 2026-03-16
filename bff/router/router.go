package router

import (
	"lianxi/bff/api"
	"lianxi/bff/handler/service/goods"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	r.POST("GoodsAdd", goods.GoodsAdd)
	r.POST("/notify/pay", api.NotifyPay)
	return r
}
