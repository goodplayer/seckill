package inventory

import (
	"context"
	"errors"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goodplayer/seckill/realworld/shared"
)

func WithholdInventory(req *shared.InventoryWithholdRequest, pool *pgxpool.Pool) error {

	fn := func() (updateCount int64, canRetry bool, rerr error) {
		recordId := uuid.Must(uuid.NewV7()).String()
		hashInventoryId := crc32.ChecksumIEEE([]byte(req.InventoryId))

		tx, err := pool.Begin(context.Background())
		if err != nil {
			rerr = err
			return
		}
		defer func() {
			_ = tx.Rollback(context.Background())
		}()

		if rows, err := tx.Query(context.Background(), `update seckill_inventory
set withholding_stock = withholding_stock + $1
where inventory_id = $2
  and pg_try_advisory_xact_lock($3)
  and withholding_stock + $1 <= total_stock returning NEW.total_stock, NEW.withholding_stock
`, req.Amount, req.InventoryId, hashInventoryId); err != nil {
			rerr = err
			return
		} else if !rows.Next() {
			defer rows.Close()
			canRetry = true
			return
		} else {
			var total int64
			var withholdingStock int64
			if err := rows.Scan(&total, &withholdingStock); err != nil {
				defer rows.Close()
				rerr = err
				return
			}
			updateCount = total - withholdingStock
			rows.Close()
		}

		if tag, err := tx.Exec(context.Background(), `
insert into seckill_inventory_order(inventory_order_record_id, inventory_id, order_id, amount, status, time_created,
                            time_updated)
values ($1, $2, $3, $4, 0, $5, $5)
`, recordId, req.InventoryId, req.OrderId, req.Amount, time.Now().UnixMilli()); err != nil {
			rerr = err
			return
		} else if tag.RowsAffected() == 0 {
			rerr = errors.New("inventory order is not inserted")
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			rerr = err
			return
		}
		return
	}

	withholdSuccess := false
	for i := 0; i < 1; i++ { // retry times
		if newStock, canRetry, err := fn(); err != nil {
			return err
		} else if canRetry {
			continue
		} else {
			var _ = newStock
			withholdSuccess = true
			break
		}
	}
	if !withholdSuccess {
		return errors.New("WithholdInventory not success due to some reason")
	}
	return nil
}
