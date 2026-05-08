package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goodplayer/seckill/realworld/shared"
)

func QueryUserById(id string, pool *pgxpool.Pool) (*shared.User, error) {
	rows, err := pool.Query(context.Background(), "SELECT user_id, username FROM seckill_user WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.UserId, &user.Username); err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, nil
	}
}
