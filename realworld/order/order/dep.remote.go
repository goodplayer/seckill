package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"resty.dev/v3"

	"github.com/goodplayer/seckill/realworld/shared"
)

type RemoteDependencyService struct {
	Client *resty.Client
	Pool   *pgxpool.Pool
}

func (r *RemoteDependencyService) WithholdInventory(req *shared.InventoryWithholdRequest) error {
	res, err := r.Client.R().
		SetContentType("application/json").
		SetBody(map[string]any{
			"inventory_id": req.InventoryId,
			"amount":       req.Amount,
			"order_id":     req.OrderId,
		}).
		Post("http://127.0.0.1:3003/inventory/withhold")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return errors.New(fmt.Sprint("withhold inventory failed:", res.Status()))
	}
	return nil
}

func (r *RemoteDependencyService) QueryUserById(id string) (*shared.User, error) {
	res, err := r.Client.R().
		SetResult(&shared.User{}).
		Get("http://127.0.0.1:3001/user/" + id)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprint("QueryUserById failed:", res.Status()))
	}

	return res.Result().(*shared.User), nil
}

func (r *RemoteDependencyService) QueryItemById(id string) (*shared.Item, error) {
	res, err := r.Client.R().
		SetResult(&shared.Item{}).
		Get("http://127.0.0.1:3002/item/" + id)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprint("QueryItemById failed:", res.Status()))
	}

	return res.Result().(*shared.Item), nil
}

func (r *RemoteDependencyService) CreateOrder(order *Order) (string, error) {
	if tag, err := r.Pool.Exec(context.Background(), `INSERT INTO seckill_order (order_id, user_id, seller_id, order_item,
                           amount, unit_price, total_price,
                           order_status, time_created, time_updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, -1, $8, $8);
`, order.OrderId, order.UserId, order.SellerId, order.ItemId, order.Amount, order.UnitPrice, order.TotalPrice, time.Now().UnixMilli()); err != nil {
		return "", err
	} else if tag.RowsAffected() != 1 {
		return "", errors.New("CreateOrder failed in db")
	}

	return order.OrderId, nil
}

func (r *RemoteDependencyService) ActivateOrder(orderId string) error {
	if tag, err := r.Pool.Exec(context.Background(), `update seckill_order 
set order_status = 1, time_updated = $2 where order_id = $1`, orderId, time.Now().UnixMilli()); err != nil {
		return err
	} else if tag.RowsAffected() != 1 {
		return errors.New("ActivateOrder failed in db")
	}
	return nil
}
