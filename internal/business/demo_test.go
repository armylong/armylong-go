package business

import (
	"context"
	"fmt"
	"testing"
)

func TestDemoBusiness_GetMessage(t *testing.T) {
	message, err := DemoBusiness.GetMessage(context.Background())
	if err != nil {
		t.Errorf("GetMessage failed: %v", err)
	}
	if message != "Hello, World!" {
		t.Errorf("GetMessage returned %v, want Hello, World!", message)
	}
	fmt.Println(message)
}
