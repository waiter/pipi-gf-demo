package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/protocol/goai"
	"pipi.com/gogf/pipi-gf-demo/internal/consts"
	"pipi.com/gogf/pipi-gf-demo/internal/controller"
	"pipi.com/gogf/pipi-gf-demo/internal/service"
	"pipi.com/gogf/pipi-gf-demo/internal/service/socket"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server of simple goframe demos",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Use(ghttp.MiddlewareHandlerResponse)
			s.BindHandler("/ws", controller.WebSocket.Upgrade)
			s.Group("/", func(group *ghttp.RouterGroup) {
				// Group middlewares.
				group.Middleware(
					service.Middleware().Ctx,
					ghttp.MiddlewareCORS,
					// service.Middleware().CORS,
				)
				// Register route handlers.
				group.Bind(
					controller.User,
				)
				// forTest
				group.ALL("/push", func(r *ghttp.Request) {
					socket.SocketManager.BroadcastMsg2Web(g.Map{
						"cmd":  "push",
						"data": "nothing",
					})
				})
				// Special handler that needs authentication.
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(service.Middleware().Auth)
					group.ALLMap(g.Map{
						"/user/profile": controller.User.Profile,
						// "/ws":           controller.WebSocket.Upgrade,
					})
				})
			})
			// Custom enhance API document.
			enhanceOpenAPIDoc(s)
			// client socket server
			controller.ClientSocket.Start()
			// Just run the server.
			s.Run()
			return nil
		},
	}
)

func enhanceOpenAPIDoc(s *ghttp.Server) {
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`

	// API description.
	openapi.Info = goai.Info{
		Title:       consts.OpenAPITitle,
		Description: consts.OpenAPIDescription,
		Contact: &goai.Contact{
			Name: "GoFrame",
			URL:  "https://goframe.org",
		},
	}
}
