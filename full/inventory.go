package full

import (
	"errors"
	"sync"

	"github.com/goodplayer/seckill/global"
)

var (
	lock       = &sync.RWMutex{}
	localCache map[int64]int64

	hotItemList      map[int64]chan *ReduceInventoryReq
	hotItemBatchSize = 10
)

func init() {
	localCache = make(map[int64]int64)

	hotItemList = make(map[int64]chan *ReduceInventoryReq)
	hotItemList[global.MIN_ITEM_ID] = make(chan *ReduceInventoryReq, 1024)
	go hotItemReduceInventoryProcessor(global.MIN_ITEM_ID, hotItemList[global.MIN_ITEM_ID])
}

func SetCacheItemQuantity(itemId, quantity int64) {
	lock.Lock()
	localCache[itemId] = quantity
	lock.Unlock()
}

func QueryInventory(itemId int64, useCache bool) (int64, error) {
	if useCache {
		lock.RLock()
		q, ok := localCache[itemId]
		lock.RUnlock()
		if ok {
			return q, nil
		}
	}

	row, err := inventory_pg_pool.Query("select quantity from item_inventory where item_id = $1 and status = 0", itemId)
	if err != nil {
		return 0, err
	}
	defer row.Close()
	if row.Next() {
		var quantity int64
		err = row.Scan(&quantity)
		if err != nil {
			return 0, err
		}
		SetCacheItemQuantity(itemId, quantity)
		return quantity, nil
	} else {
		return 0, errors.New("item not found or item status abnormal.")
	}
}

func ReduceInventory(itemId, quantity int64) (int64, error) {
	c, ok := hotItemList[itemId]
	if ok {
		// batch
		req := &ReduceInventoryReq{
			Quantity: quantity,
			Result:   make(chan error, 1),
		}
		c <- req
		err := <-req.Result
		return 0, err
	} else {
		return reduceInventoryInternal(itemId, quantity)
	}
}

func reduceInventoryInternal(itemId, quantity int64) (int64, error) {
	r, err := inventory_pg_pool.Query("update item_inventory set quantity = quantity - $1 where item_id = $2 and status = 0 and quantity >= $3 returning quantity", quantity, itemId, quantity)
	if err != nil {
		return 0, err
	} else {
		defer r.Close()
		if r.Next() {
			var newQuantity int64
			err = r.Scan(&newQuantity)
			if err != nil {
				return 0, err
			}
			err = UpdateQuantityCache(itemId, newQuantity)
			if err != nil {
				return 0, err
			}
			return newQuantity, nil
		} else {
			QueryInventory(itemId, false)
			return 0, errors.New("no item inventory reduced.")
		}
	}
}

func UpdateQuantityCache(itemId, quantity int64) error {
	//TODO quantity cache
	SetCacheItemQuantity(itemId, quantity)
	return nil
}

type ReduceInventoryReq struct {
	Quantity int64
	Result   chan error
}

func hotItemReduceInventoryProcessor(itemId int64, rc chan *ReduceInventoryReq) {
	var i = 0
	reqList := make([]*ReduceInventoryReq, hotItemBatchSize)
	for {
		select {
		case item := <-rc:
			reqList[i] = item
			i++
		default:
			if i > 0 {
				BatchReduceInventory(itemId, reqList[:i])
				i = 0
			} else {
				item := <-rc
				reqList[i] = item
				i++
			}
		}
		if i >= hotItemBatchSize {
			BatchReduceInventory(itemId, reqList[:i])
			i = 0
		}
	}
}

func BatchReduceInventory(itemId int64, reqList []*ReduceInventoryReq) {
	var all int64 = 0
	for _, v := range reqList {
		all += v.Quantity
	}

	_, err := reduceInventoryInternal(itemId, all)
	if err != nil {
		// run each reduce
		for _, v := range reqList {
			_, err := reduceInventoryInternal(itemId, v.Quantity)
			v.Result <- err
		}
	} else {
		for _, v := range reqList {
			v.Result <- nil
		}
	}
}
