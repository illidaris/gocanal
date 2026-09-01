package server

import (
	"context"
	"fmt"
	"gocanal/internal/canal"
	canalEx "gocanal/pkg/canal"
	"strings"

	"github.com/illidaris/aphrodite/ginhandle/middleware"
	"github.com/illidaris/aphroditecli/pkg/log"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/ginhandle"

	ginEx "github.com/illidaris/gin"
)

func Run() {
	ctx := context.Background()
	engine := ginhandle.NewGin(
		// ginhandle.WithMode(config.GetString("app.mode")),
		ginhandle.WithParamMiddleware(true,
			middleware.WithRequestContentLengthMax(1048576),
			middleware.WithResponseContentLengthMax(1048576),
			middleware.WithAfterFunc(func(subCtx context.Context, info *middleware.APIInfo) {
				if info != nil && info.Cost > (time.Second*10).Milliseconds() {
					log.Warn(subCtx, "[SLOW]请求接口：%s, 请求方法：%s, 请求耗时：%dms", info.Path, info.Method, info.Cost)
				}
			}),
		),
		ginhandle.WithMetricReqCntURLLabelMappingFn(func(c *gin.Context) string {
			url := c.Request.URL.Path
			if c.Request.Method == "OPTIONS" {
				return "OPTIONS"
			}
			for _, param := range c.Params {
				// 聚合ID路由
				if param.Key == "id" {
					url = strings.Replace(url, param.Value, ":"+param.Key, 1)
					break
				}
			}
			return url
		}),
		ginhandle.WithCollectors(canalEx.ReqNum, canalEx.DeleteNum, canalEx.IndexNum),
		ginhandle.WithGinInnerHandle(func(e *gin.Engine) {}),
	)
	canal.Go(ctx)

	port := 8080
	lsn := fmt.Sprintf(":%d", port)
	// run
	ginEx.GracefulRunWithAop(ctx, engine, lsn, time.Second*5, func(_ int) {}, func() {})
}
