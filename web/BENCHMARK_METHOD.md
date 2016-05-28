### hot item benchmark

```
batch size : hotItemBatchSize = 10

./wrk --timeout 10s -c 1000 -t 50 -d 10s "http://localhost:7649/reducing?item_id=10000000000&user_id=100000"
```

### normal item

```
./wrk --timeout 10s -c 1000 -t 50 -d 10s "http://localhost:7649/reducing?item_id=10000000001&user_id=100000"
```

### small quantity inventory item

```
inventory size = 10

./wrk --timeout 10s -c 1000 -t 50 -d 10s "http://localhost:7649/reducing?item_id=2000000000&user_id=100000"
```

### two db item

```
./wrk --timeout 10s -c 1000 -t 50 -d 10s "http://localhost:7649/reducing?item_id=3000000000&user_id=100000"
```