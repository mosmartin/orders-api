package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mosmartin/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

type FindAllPage struct {
	Size   uint
	Offset uint
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func orderIDKey(orderID uint64) string {
	return fmt.Sprintf("order:%d", orderID)
}

func (r *RedisRepository) List(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to list orders: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to list orders: %w", err)
	}

	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		var order model.Order

		err := json.Unmarshal([]byte(x.(string)), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order: %w", err)
		}

		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}

func (r *RedisRepository) Create(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		slog.Error("failed to encode order", err)
	}

	key := orderIDKey(order.OrderID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if res.Err() != nil {
		txn.Discard()

		return fmt.Errorf("failed to insert order: %w", res.Err())
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()

		return fmt.Errorf("failed to insert order: %w", err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to exec order: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepository) GetByID(ctx context.Context, orderID uint64) (model.Order, error) {
	key := orderIDKey(orderID)

	data, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order model.Order

	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order: %w", err)
	}

	return order, nil
}

func (r *RedisRepository) UpdateByID(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		slog.Error("failed to encode order", err)
	}

	key := orderIDKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func (r *RedisRepository) DeleteByID(ctx context.Context, orderID uint64) error {
	key := orderIDKey(orderID)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()

		return ErrNotExist
	} else if err != nil {
		txn.Discard()

		return fmt.Errorf("failed to delete order: %w", err)
	}

	err = txn.SRem(ctx, "orders", key).Err()
	if err != nil {
		txn.Discard()

		return fmt.Errorf("failed to delete order: %w", err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to exec order: %w", err)
	}

	return nil
}
