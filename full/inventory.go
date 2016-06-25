package full

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"sync"

	"github.com/goodplayer/seckill/global"
)

var (
	lock       = &sync.RWMutex{}
	localCache map[int64]int64

	hotItemList      map[int64]chan *ReduceInventoryReq
	hotItemBatchSize = 10

	multiDbItem      map[int64]bool
	multiDbItemCache map[int64]MultiDbItemCacheLine
	multiDbitemLock  = &sync.Mutex{}
)

var (
	queryInventoryKey          string
	reduceAndQueryInventoryKey string
)

func prepareInventorySql() {
	queryInventory := "select quantity from item_inventory where item_id = $1 and status = 0"
	queryInventoryMd5 := md5.Sum([]byte(queryInventory))
	queryInventoryKey = "queryInventory_" + hex.EncodeToString(queryInventoryMd5[:])
	_, err := inventory_pg_pool.Prepare(queryInventoryKey, queryInventory)
	if err != nil {
		log.Fatalln("prepare sql pg1 - queryInventory error", err)
	} else {
		log.Println("prepare sql pg1 - queryInventory - key:", queryInventoryKey, "sql:", queryInventory)
	}
	_, err = inventory_pg2_pool.Prepare(queryInventoryKey, queryInventory)
	if err != nil {
		log.Fatalln("prepare sql pg2 - queryInventory error", err)
	} else {
		log.Println("prepare sql pg2 - queryInventory - key:", queryInventoryKey, "sql:", queryInventory)
	}
	reduceAndQueryInventory := "update item_inventory set quantity = quantity - $1 where item_id = $2 and status = 0 and quantity >= $3 and pg_try_advisory_xact_lock($2) returning quantity"
	reduceAndQueryInventoryMd5 := md5.Sum([]byte(reduceAndQueryInventory))
	reduceAndQueryInventoryKey = "reduceAndQuery_" + hex.EncodeToString(reduceAndQueryInventoryMd5[:])
	_, err = inventory_pg_pool.Prepare(reduceAndQueryInventoryKey, reduceAndQueryInventory)
	if err != nil {
		log.Fatalln("prepare sql pg1 - reduceAndQueryInventory error", err)
	} else {
		log.Println("prepare sql pg1 - reduceAndQueryInventory -- key:", reduceAndQueryInventoryKey, "sql:", reduceAndQueryInventory)
	}
	_, err = inventory_pg_pool.Prepare(reduceAndQueryInventoryKey, reduceAndQueryInventory)
	if err != nil {
		log.Fatalln("prepare sql pg2 - reduceAndQueryInventory error", err)
	} else {
		log.Println("prepare sql pg2 - reduceAndQueryInventory -- key:", reduceAndQueryInventoryKey, "sql:", reduceAndQueryInventory)
	}
}

// currently only support 2 dbs
type MultiDbItemCacheLine struct {
	Count [2]int64
	Total int64
}

func init() {
	localCache = make(map[int64]int64)

	multiDbItem = make(map[int64]bool)
	multiDbItem[3000000000] = true
	multiDbItemCache = make(map[int64]MultiDbItemCacheLine)

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
	if v, ok := multiDbItem[itemId]; ok && v {
		return queryMultiInventory(itemId, useCache)
	}

	if useCache {
		lock.RLock()
		q, ok := localCache[itemId]
		lock.RUnlock()
		if ok {
			return q, nil
		}
	}

	row, err := inventory_pg_pool.Query(queryInventoryKey, itemId)
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
	if v, ok := multiDbItem[itemId]; ok && v {
		_, err := reduceMultiInventory(itemId, quantity)
		// query different Pool, must be out of Row. Otherwise may deadlock. especially when use defer
		cnt, err2 := queryMultiInventory(itemId, false)
		if err != nil {
			return 0, err
		}
		if err2 != nil {
			return 0, err2
		}
		return cnt, nil
	}

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

func queryMultiInventory(itemId int64, useCache bool) (int64, error) {
	if useCache {
		multiDbitemLock.Lock()
		line, ok := multiDbItemCache[itemId]
		multiDbitemLock.Unlock()
		if ok {
			return line.Total, nil
		}
	}

	c1, err := loadDb1Count(itemId)
	if err != nil {
		return 0, err
	}
	c2, err := loadDb2Count(itemId)
	if err != nil {
		return 0, err
	}

	line := MultiDbItemCacheLine{
		Count: [2]int64{c1, c2},
		Total: c1 + c2,
	}

	multiDbitemLock.Lock()
	multiDbItemCache[itemId] = line
	multiDbitemLock.Unlock()
	return line.Total, nil
}

func loadDb1Count(itemId int64) (int64, error) {
	row, err := inventory_pg_pool.Query(queryInventoryKey, itemId)
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
		return quantity, nil
	} else {
		return 0, errors.New("item not found or item status abnormal in db1. " + row.Err().Error())
	}
}

func loadDb2Count(itemId int64) (int64, error) {
	row, err := inventory_pg2_pool.Query(queryInventoryKey, itemId)
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
		return quantity, nil
	} else {
		return 0, errors.New("item not found or item status abnormal in db2. " + row.Err().Error())
	}
}

func reduceMultiInventory(itemId, quantity int64) (bool, error) {
	i := rand.Int31n(2)
	if i == 0 {
		row, err := inventory_pg_pool.Query(reduceAndQueryInventoryKey, quantity, itemId, quantity)
		if err != nil {
			return false, err
		} else {
			defer row.Close()
			if row.Next() {
				return true, nil
			} else {
				return false, errors.New("no item inventory reduced.")
			}
		}
	} else if i == 1 {
		row, err := inventory_pg2_pool.Query(reduceAndQueryInventoryKey, quantity, itemId, quantity)
		if err != nil {
			return false, err
		} else {
			defer row.Close()
			if row.Next() {
				return true, nil
			} else {
				return false, errors.New("no item inventory reduced.")
			}
		}
	} else {
		return false, errors.New("error reduce multi-db item. random number is not 0 or 1")
	}
}

func reduceInventoryInternal(itemId, quantity int64) (int64, error) {
	r, err := inventory_pg_pool.Query(reduceAndQueryInventoryKey, quantity, itemId, quantity)
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
	// quantity cache
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
