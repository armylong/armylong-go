package controllers

import (
	"errors"
	"net/http"

	"github.com/armylong/armylong-go/internal/controllers/user"

	"github.com/armylong/go-library/service/longgin"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(engine *gin.Engine) {

	engine.LoadHTMLGlob(`./templates/*.gohtml`)

	// 一个简单的示例，访问/，直接输出一个文本
	engine.Any(`/`, homepage)

	// 404 页面
	engine.NoRoute(func(ctx *gin.Context) {
		// 返回404
		ctx.JSON(http.StatusNotFound, longgin.ErrorWithContext(ctx, errors.New("not found"), longgin.ErrorNotFound))
	})

	// 用户相关接口
	userRoot := engine.Group("/user")
	longgin.RegisterJsonController(userRoot.Group("/demo"), &user.DemoController{})

}
