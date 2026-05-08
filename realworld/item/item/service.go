package item

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goodplayer/seckill/realworld/shared"
)

func QueryItemById(id string, pool *pgxpool.Pool) (*shared.Item, error) {
	rows, err := pool.Query(context.Background(), "SELECT item_id, seller_id, item_name, description, unit_price, inventory_id FROM seckill_item WHERE item_id = $1 and status = 0", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		row := shared.Item{}
		if err := rows.Scan(&row.ItemId, &row.SellerId, &row.ItemName, &row.Description, &row.UnitPrice, &row.InventoryId); err != nil {
			return nil, err
		}
		return &row, nil
	} else {
		return nil, nil
	}
}
