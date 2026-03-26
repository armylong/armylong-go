package user

import (
	userCs "github.com/armylong/armylong-go/internal/cs/user"
	"github.com/gin-gonic/gin"
)

type DemoController struct {
}

func (u *DemoController) ActionHello(c *gin.Context) (res userCs.DemoResponse, err error) {
	res = userCs.DemoResponse{Message: "Hello, World!"}
	return res, nil
}
