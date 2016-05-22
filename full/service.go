package full

import "errors"

type CreateOrderReq struct {
	UserId      int64
	ItemId      int64
	BuyQuantity int64
}

type CreateOrderResp struct {
	OrderId int64
}

func CreateOrder(req *CreateOrderReq) (*CreateOrderResp, error) {
	//TODO 1. ratelimit

	// 2. query inventory
	quantity, err := QueryInventory(req.ItemId)
	if err != nil {
		return nil, err
	}
	if quantity < req.BuyQuantity {
		return nil, errors.New("inventory not enough")
	}

	// 3. generate order id
	orderId, err := generateOrderId(req.UserId)
	if err != nil {
		return nil, err
	}

	// 3. reduce inventory
	err = ReduceInventory(req.ItemId, req.BuyQuantity)
	if err != nil {
		return nil, err
	}

	// 4. create order
	err = SaveOrder(orderId, req.UserId, req.ItemId)
	if err != nil {
		//TODO 4.5 add back inventory
		return nil, err
	}

	resp := &CreateOrderResp{
		OrderId: orderId,
	}

	return resp, nil
}
