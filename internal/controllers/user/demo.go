package user

import (
	"github.com/gin-gonic/gin"
)

type DemoController struct {
}

type response struct {
	Message string `json:"message"`
}

func (u *DemoController) ActionHello(c *gin.Context) (res response, err error) {
	res = response{Message: "Hello, World!"}
	return res, nil
}
