package user

import (
	"github.com/armylong/armylong-go/internal/business"
	userCs "github.com/armylong/armylong-go/internal/cs/user"
	"github.com/gin-gonic/gin"
)

type DemoController struct {
}

func (c *DemoController) ActionHello(ctx *gin.Context) (res *userCs.DemoResponse, err error) {
	message, err := business.DemoBusiness.GetMessage(ctx)
	if err != nil || message == "" {
		return res, err
	}

	return &userCs.DemoResponse{
		Message: message,
	}, nil
}
