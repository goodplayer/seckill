package order

import (
	"github.com/goodplayer/seckill/realworld/shared"
)

type LocalDependencyService struct {
}

func (l *LocalDependencyService) WithholdInventory(req *shared.InventoryWithholdRequest) error {
	return nil
}

func (l *LocalDependencyService) QueryUserById(id string) (*shared.User, error) {
	return &shared.User{
		UserId:   id,
		Username: "testuser",
	}, nil
}

func (l *LocalDependencyService) QueryItemById(id string) (*shared.Item, error) {
	return &shared.Item{
		ItemId:      id,
		SellerId:    "1111",
		ItemName:    "itemdemo",
		Description: "This is demo.",
		UnitPrice:   2,
		InventoryId: "2222",
	}, nil
}

func (l *LocalDependencyService) CreateOrder(order *Order) (string, error) {
	return order.OrderId, nil
}

func (l *LocalDependencyService) ActivateOrder(orderId string) error {
	return nil
}
