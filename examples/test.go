package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

// gin.ctx不建议在协程内使用
// see:https://jinfeijie.cn/post-777
func main() {
	g := gin.Default()
	g.GET("/test", test)
	g.Run(":8080")
}

func test(ctx *gin.Context) {
	_uuid := uuid.New().String()
	ctx.Set("uuid", _uuid)
	go func(c *gin.Context, u string) {
		time.Sleep(time.Second)
		cu := c.GetString("uuid")
		if cu == u {
			log.Printf("一致 %p, %s\n", c, u)
		} else {
			log.Printf("不一致 %p, uuid = %s, cuuid = %s\n", c, u, cu)
		}
	}(ctx, _uuid)

	ctx.JSON(200, _uuid)
	return
}
