package user

import (
	"github.com/gin-gonic/gin"
)

type DemoController struct {
}

func (u *DemoController) ActionHello(c *gin.Context) (res string, err error) {

	return "Hello, World!", nil
}
