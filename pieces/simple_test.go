package pieces

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"xorm.io/xorm"
)

func TestBasicInsert(t *testing.T) {
	engine, err := xorm.NewEngine("pgx", "postgres://admin:admin@10.11.0.7:5432/inventory")
	if err != nil {
		t.Error(err)
	}
	engine.SetMaxOpenConns(150)
	engine.ShowSQL(false)

	const workerCnt = 100
	const taskCntPerWorker = 100
	barrier := make(chan struct{})
	wg := new(sync.WaitGroup)
	wg.Add(workerCnt)
	for i := 0; i < workerCnt; i++ {
		go func() {
			<-barrier
			defer wg.Done()

			for i := 0; i < taskCntPerWorker; i++ {
				if err := deduct(engine.NewSession(), 1, 1); err != nil {
					t.Error(err)
				}
			}
		}()
	}

	log.Println("starting....")
	start := time.Now()
	close(barrier)
	wg.Wait()
	fmt.Println(time.Since(start))
}

func deduct(sess *xorm.Session, itemId, quantity int64) error {
	now := time.Now().UnixMilli()
	defer func(sess *xorm.Session) {
		err := sess.Close()
		if err != nil {
			log.Println("close session failed:", err)
		}
	}(sess)
	if err := sess.Begin(); err != nil {
		return err
	}
	defer func(sess *xorm.Session) {
		err := sess.Rollback()
		if err != nil {
			log.Println("rollback failed:", err)
		}
	}(sess)

	if r, err := sess.QuerySliceString("update inventory set quantity = quantity - ?, time_updated = ? where item_id = ? and quantity >= ? and pg_try_advisory_xact_lock(?) returning item_id, quantity",
		quantity, now, itemId, quantity, itemId); err != nil {
		return err
	} else {
		var _ = r
		//for _, col := range r {
		//	log.Println(col)
		//}
	}

	if err := sess.Commit(); err != nil {
		return err
	}
	return nil
}
