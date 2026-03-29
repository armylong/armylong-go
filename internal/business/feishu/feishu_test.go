package feishu

import (
	"context"
	"fmt"
	"testing"

	feishuCs "github.com/armylong/armylong-go/internal/cs/feishu"
	feishuLibrary "github.com/armylong/go-library/service/feishu"
)

var ctx = context.Background()

func TestGetAuthorizationHeader(t *testing.T) {
	BitableBusiness.GetBitable(ctx, &feishuCs.BitableRequest{
		AppToken: "6BHkHf8DJAbdA608F9azx89x2xK9LHwd",
	})
}

func TestInitUserAccessToken(t *testing.T) {
	code := "bDKkDf7Kyy88Ab18Hfz7KyI20ybzbH4K"
	redirectUri := "https://olfrzjwptnle.ap-northeast-1.clawcloudrun.com"
	userAccessTokenHeader := feishuLibrary.GetUserAccessTokenHeader(code, redirectUri)
	fmt.Println(userAccessTokenHeader)
}

func TestRefreshUserAccessToken(t *testing.T) {
	userAccessTokenHeader := feishuLibrary.GetUserAccessTokenHeader(``, ``)
	fmt.Println(userAccessTokenHeader)
}
