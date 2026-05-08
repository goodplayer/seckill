package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/goodplayer/seckill/realworld/shared"
)

type OrderService struct {
	DependencyService DependencyService
}

func (o *OrderService) PlaceOrder(req *shared.PlaceOrderRequest) (string, error) {
	user, err := o.DependencyService.QueryUserById(req.UserId)
	if err != nil {
		return "", err
	} else if user == nil {
		return "", errors.New("user not found")
	}
	item, err := o.DependencyService.QueryItemById(req.ItemId)
	if err != nil {
		return "", err
	} else if item == nil {
		return "", errors.New("item not found")
	}

	order, err := o.performBusinesses(req, user, item)
	if err != nil {
		return "", err
	}

	// save order
	orderId, err := o.DependencyService.CreateOrder(order)
	if err != nil {
		return "", err
	}

	// deduct resources
	if err := o.deductResources(order, item); err != nil {
		return "", err
	}

	// activate order
	if err := o.DependencyService.ActivateOrder(orderId); err != nil {
		return "", err
	}
	return orderId, nil
}

func (o *OrderService) performBusinesses(req *shared.PlaceOrderRequest, user *shared.User, item *shared.Item) (*Order, error) {
	order := new(Order)

	// businesses: discount calculation, logistic calculation, limitation calculation, security check

	// price calculation
	priceFn := func() {
		order.TotalPrice = item.UnitPrice * req.Amount
		order.UnitPrice = item.UnitPrice
		order.Amount = req.Amount
	}
	priceFn()

	// order assembly
	orderFn := func() {
		order.UserId = req.UserId
		order.SellerId = item.SellerId
		order.ItemId = req.ItemId
		order.OrderId = uuid.Must(uuid.NewV7()).String()
	}
	orderFn()

	return order, nil
}

func (o *OrderService) deductResources(order *Order, item *shared.Item) error {
	// withhold inventory
	req := item.ToInventoryWithholdRequest(order.Amount, order.OrderId)
	if err := o.DependencyService.WithholdInventory(req); err != nil {
		return err
	}
	return nil
}

type DependencyService interface {
	QueryUserById(id string) (*shared.User, error)
	QueryItemById(id string) (*shared.Item, error)

	WithholdInventory(req *shared.InventoryWithholdRequest) error

	CreateOrder(order *Order) (string, error)
	ActivateOrder(orderId string) error
}

type Order struct {
	OrderId  string
	UserId   string
	SellerId string
	ItemId   string

	Amount     int64
	UnitPrice  int64
	TotalPrice int64
}
