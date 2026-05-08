package shared

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type Item struct {
	ItemId      string `json:"item_id"`
	SellerId    string `json:"seller_id"`
	ItemName    string `json:"item_name"`
	Description string `json:"description"`
	UnitPrice   int64  `json:"unit_price"` // may not exceed the max allowed value in json
	InventoryId string `json:"inventory_id"`
}

func (i *Item) ToInventoryWithholdRequest(amount int64, orderId string) *InventoryWithholdRequest {
	return &InventoryWithholdRequest{
		Amount:      amount,
		InventoryId: i.InventoryId,
		OrderId:     orderId,
	}
}

type InventoryWithholdRequest struct {
	InventoryId string `json:"inventory_id"`
	Amount      int64  `json:"amount"`
	OrderId     string `json:"order_id"`
}

type PlaceOrderRequest struct {
	UserId string `json:"user_id"`
	ItemId string `json:"item_id"`
	Amount int64  `json:"amount"`
}
