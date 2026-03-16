package goods

import (
	"lianxi/bff/dasic/config"
	"lianxi/bff/handler/request"
	__ "lianxi/srv/proto/goods"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoodsAdd(c *gin.Context) {
	var form request.Goods
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 400,
		})
		return
	}
	r, err := config.GoodsClient.GoodsAdd(c, &__.GoodsAddReq{
		Name:  form.Name,
		Price: uint32(form.Price),
		Num:   form.Num,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":  r.Msg,
		"code": r.Code,
	})
	return
}
