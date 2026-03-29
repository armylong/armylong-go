package feishu

import (
	"context"
	"fmt"

	feishuCs "github.com/armylong/armylong-go/internal/cs/feishu"
	feishuLibrary "github.com/armylong/go-library/service/feishu"
)

type bitableBusiness struct{}

var BitableBusiness = &bitableBusiness{}

func (b *bitableBusiness) GetBitable(ctx context.Context, req *feishuCs.BitableRequest) ([]byte, error) {
	// code := "fzLsCc5CBwefA320Gbz5yA3cc053KK9f"
	// redirectUri := "https://olfrzjwptnle.ap-northeast-1.clawcloudrun.com"
	// userAccessTokenHeader := feishuLibrary.GetUserAccessTokenHeader(code, redirectUri)
	userAccessTokenHeader := feishuLibrary.GetUserAccessTokenHeader(``, ``)
	fmt.Println(userAccessTokenHeader)

	// res, err := httpx.GetWithHeader(fmt.Sprintf(feishuConfig.FeishuApiBitableUrl, req.AppToken), map[string]string{
	// 	`Authorization`: userAccessTokenHeader,
	// 	`Content-Type`:  `application/json`,
	// })
	// fmt.Println(string(res), err)

	return nil, nil
	// ZFszben8BaPhvPscIbLcmKsZnYB
}
