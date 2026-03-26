package business

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/armylong/armylong-go/internal/common/webcache"
)

const (
	YangfenBalanceKey     = "yangfen:balance:"
	YangfenTransactionKey = "yangfen:transactions:"
	YangfenExpireKey      = "yangfen:expire:"
)

type yangfenBusiness struct{}

var YangfenBusiness = &yangfenBusiness{}

func (b *yangfenBusiness) getBalanceKey(uid string) string {
	return YangfenBalanceKey + uid
}

func (b *yangfenBusiness) getTransactionKey(uid string) string {
	return YangfenTransactionKey + uid
}

func (b *yangfenBusiness) getExpireKey(uid string) string {
	return YangfenExpireKey + uid
}

func (b *yangfenBusiness) GetBalance(ctx context.Context, uid string) (int, error) {
	balanceStr, err := webcache.RedisClient.Get(ctx, b.getBalanceKey(uid)).Result()
	if err != nil {
		return 0, nil
	}
	balance, _ := strconv.Atoi(balanceStr)
	return balance, nil
}

func (b *yangfenBusiness) checkAndClearExpired(ctx context.Context, uid string) error {
	expireKey := b.getExpireKey(uid)
	expireStr, err := webcache.RedisClient.Get(ctx, expireKey).Result()
	if err != nil {
		return nil
	}
	expireTime, _ := strconv.ParseInt(expireStr, 10, 64)
	if time.Now().Unix() > expireTime {
		webcache.RedisClient.Set(ctx, b.getBalanceKey(uid), "0", 0)
		webcache.RedisClient.Del(ctx, expireKey)
	}
	return nil
}

func (b *yangfenBusiness) Recharge(ctx context.Context, uid string, amount int, expireSec int64) error {
	if amount <= 0 {
		return fmt.Errorf("充值金额必须大于0")
	}

	b.checkAndClearExpired(ctx, uid)

	balance, _ := b.GetBalance(ctx, uid)
	newBalance := balance + amount

	webcache.RedisClient.Set(ctx, b.getBalanceKey(uid), strconv.Itoa(newBalance), 0)

	expireTime := time.Now().Add(time.Duration(expireSec) * time.Second).Unix()
	webcache.RedisClient.Set(ctx, b.getExpireKey(uid), strconv.FormatInt(expireTime, 10), 0)

	b.addTransaction(ctx, uid, "recharge", amount, newBalance, "充值")
	return nil
}

func (b *yangfenBusiness) Consume(ctx context.Context, uid string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("消费金额必须大于0")
	}

	b.checkAndClearExpired(ctx, uid)

	balance, _ := b.GetBalance(ctx, uid)
	if balance < amount {
		return fmt.Errorf("余额不足")
	}

	newBalance := balance - amount
	webcache.RedisClient.Set(ctx, b.getBalanceKey(uid), strconv.Itoa(newBalance), 0)

	bonusRate := 1
	if amount >= 100 {
		bonusRate = 2
	}
	_ = bonusRate

	b.addTransaction(ctx, uid, "consume", amount, newBalance, fmt.Sprintf("消费%d积分", amount))
	return nil
}

func (b *yangfenBusiness) Transfer(ctx context.Context, fromUid, toUid string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("转账金额必须大于0")
	}
	if fromUid == toUid {
		return fmt.Errorf("不能转给自己")
	}

	b.checkAndClearExpired(ctx, fromUid)
	b.checkAndClearExpired(ctx, toUid)

	fromBalance, _ := b.GetBalance(ctx, fromUid)
	if fromBalance < amount {
		return fmt.Errorf("余额不足")
	}

	newFromBalance := fromBalance - amount
	webcache.RedisClient.Set(ctx, b.getBalanceKey(fromUid), strconv.Itoa(newFromBalance), 0)

	toBalance, _ := b.GetBalance(ctx, toUid)
	newToBalance := toBalance + amount
	webcache.RedisClient.Set(ctx, b.getBalanceKey(toUid), strconv.Itoa(newToBalance), 0)

	b.addTransaction(ctx, fromUid, "transfer_out", amount, newFromBalance, fmt.Sprintf("转出给%s", toUid))
	b.addTransaction(ctx, toUid, "transfer_in", amount, newToBalance, fmt.Sprintf("从%s转入", fromUid))
	return nil
}

func (b *yangfenBusiness) Refund(ctx context.Context, uid string, transactionId string) error {
	transaction, err := b.getTransaction(ctx, uid, transactionId)
	if err != nil {
		return fmt.Errorf("交易记录不存在")
	}

	txType, _ := transaction["type"].(string)
	if txType != "consume" {
		return fmt.Errorf("只能退款消费记录")
	}

	balance, _ := b.GetBalance(ctx, uid)

	refundAmount := int(transaction["amount"].(float64))
	newBalance := balance + refundAmount
	webcache.RedisClient.Set(ctx, b.getBalanceKey(uid), strconv.Itoa(newBalance), 0)

	b.addTransaction(ctx, uid, "refund", refundAmount, newBalance, fmt.Sprintf("退款-交易号:%s", transactionId))
	return nil
}

func (b *yangfenBusiness) GetTransactions(ctx context.Context, uid string) ([]map[string]any, error) {
	key := b.getTransactionKey(uid)
	data, err := webcache.RedisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	transactions := make([]map[string]any, 0, len(data))
	for _, item := range data {
		var t map[string]any
		json.Unmarshal([]byte(item), &t)
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (b *yangfenBusiness) addTransaction(ctx context.Context, uid string, txType string, amount int, balance int, desc string) {
	transaction := map[string]any{
		"id":          fmt.Sprintf("TX%d", time.Now().UnixNano()),
		"uid":         uid,
		"type":        txType,
		"amount":      amount,
		"balance":     balance,
		"description": desc,
		"createdAt":   time.Now().Unix(),
	}

	data, _ := json.Marshal(transaction)
	key := b.getTransactionKey(uid)
	webcache.RedisClient.LPush(ctx, key, string(data))
	webcache.RedisClient.LTrim(ctx, key, 0, 99)
}

func (b *yangfenBusiness) getTransaction(ctx context.Context, uid string, transactionId string) (map[string]any, error) {
	transactions, _ := b.GetTransactions(ctx, uid)
	for _, t := range transactions {
		if t["id"] == transactionId {
			return t, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (b *yangfenBusiness) ClearData(ctx context.Context, uid string) error {
	webcache.RedisClient.Del(ctx, b.getBalanceKey(uid))
	webcache.RedisClient.Del(ctx, b.getTransactionKey(uid))
	webcache.RedisClient.Del(ctx, b.getExpireKey(uid))
	return nil
}
