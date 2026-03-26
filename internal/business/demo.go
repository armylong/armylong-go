package business

import (
	"context"
)

type demoBusiness struct{}

var DemoBusiness = &demoBusiness{}

func (b *demoBusiness) GetMessage(ctx context.Context) (res string, err error) {
	return "Hello, World!", nil
}
