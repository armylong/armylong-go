package yangfen

import (
	"errors"

	"github.com/armylong/armylong-go/internal/business"
	"github.com/armylong/armylong-go/internal/cs/yangfen"
	"github.com/gin-gonic/gin"
)

type YangfenController struct {
}

func (c *YangfenController) checkUid(uid string) error {
	if uid == "" {
		return errors.New("uid不能为空")
	}
	return nil
}

func (c *YangfenController) ActionGetBalance(ctx *gin.Context, req *yangfen.BaseRequest) (*yangfen.BalanceResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	balance, err := business.YangfenBusiness.GetBalance(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &yangfen.BalanceResponse{
		Uid:     req.Uid,
		Balance: balance,
	}, nil
}

func (c *YangfenController) ActionRecharge(ctx *gin.Context, req *yangfen.RechargeRequest) (*yangfen.BalanceResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	if req.Amount <= 0 {
		return nil, errors.New("充值金额必须大于0")
	}

	err := business.YangfenBusiness.Recharge(ctx, req.Uid, req.Amount, req.ExpireSec)
	if err != nil {
		return nil, err
	}

	balance, _ := business.YangfenBusiness.GetBalance(ctx, req.Uid)

	return &yangfen.BalanceResponse{
		Uid:     req.Uid,
		Balance: balance,
	}, nil
}

func (c *YangfenController) ActionConsume(ctx *gin.Context, req *yangfen.ConsumeRequest) (*yangfen.BalanceResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	if req.Amount <= 0 {
		return nil, errors.New("消费金额必须大于0")
	}

	err := business.YangfenBusiness.Consume(ctx, req.Uid, req.Amount)
	if err != nil {
		return nil, err
	}

	balance, _ := business.YangfenBusiness.GetBalance(ctx, req.Uid)

	return &yangfen.BalanceResponse{
		Uid:     req.Uid,
		Balance: balance,
	}, nil
}

func (c *YangfenController) ActionTransfer(ctx *gin.Context, req *yangfen.TransferRequest) (*yangfen.CommonResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	if req.ToUid == "" {
		return nil, errors.New("目标用户不能为空")
	}

	if req.Amount <= 0 {
		return nil, errors.New("转账金额必须大于0")
	}

	err := business.YangfenBusiness.Transfer(ctx, req.Uid, req.ToUid, req.Amount)
	if err != nil {
		return nil, err
	}

	fromBalance, _ := business.YangfenBusiness.GetBalance(ctx, req.Uid)
	toBalance, _ := business.YangfenBusiness.GetBalance(ctx, req.ToUid)

	return &yangfen.CommonResponse{
		Success: true,
		Message: "转账成功",
		Data: map[string]any{
			"fromUid":     req.Uid,
			"fromBalance": fromBalance,
			"toUid":       req.ToUid,
			"toBalance":   toBalance,
		},
	}, nil
}

func (c *YangfenController) ActionRefund(ctx *gin.Context, req *yangfen.RefundRequest) (*yangfen.BalanceResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	if req.TransactionId == "" {
		return nil, errors.New("交易号不能为空")
	}

	err := business.YangfenBusiness.Refund(ctx, req.Uid, req.TransactionId)
	if err != nil {
		return nil, err
	}

	balance, _ := business.YangfenBusiness.GetBalance(ctx, req.Uid)

	return &yangfen.BalanceResponse{
		Uid:     req.Uid,
		Balance: balance,
	}, nil
}

func (c *YangfenController) ActionGetTransactions(ctx *gin.Context, req *yangfen.BaseRequest) (*yangfen.TransactionListResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	transactions, err := business.YangfenBusiness.GetTransactions(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &yangfen.TransactionListResponse{
		List:  convertTransactions(transactions),
		Total: len(transactions),
	}, nil
}

func convertTransactions(transactions []map[string]any) []yangfen.TransactionRecord {
	result := make([]yangfen.TransactionRecord, 0, len(transactions))
	for _, t := range transactions {
		record := yangfen.TransactionRecord{}
		if id, ok := t["id"].(string); ok {
			record.Id = id
		}
		if uid, ok := t["uid"].(string); ok {
			record.Uid = uid
		}
		if txType, ok := t["type"].(string); ok {
			record.Type = txType
		}
		if amount, ok := t["amount"].(float64); ok {
			record.Amount = int(amount)
		}
		if balance, ok := t["balance"].(float64); ok {
			record.Balance = int(balance)
		}
		if desc, ok := t["description"].(string); ok {
			record.Description = desc
		}
		if createdAt, ok := t["createdAt"].(float64); ok {
			record.CreatedAt = int64(createdAt)
		}
		result = append(result, record)
	}
	return result
}

func (c *YangfenController) ActionClearData(ctx *gin.Context, req *yangfen.BaseRequest) (*yangfen.CommonResponse, error) {
	if err := c.checkUid(req.Uid); err != nil {
		return nil, err
	}

	err := business.YangfenBusiness.ClearData(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &yangfen.CommonResponse{
		Success: true,
		Message: "清除成功",
	}, nil
}
