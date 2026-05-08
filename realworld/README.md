# RealWorld

## Features

* [x] Four core services
    * User
    * Item
    * Inventory
    * Order
* [ ] RealWorld Business
    * [x] PlaceOrder - withhold inventory
    * [ ] Paid - deduct
    * [ ] Stock replenishment
    * [ ] Close order
    * [ ] Close paid order
    * [ ] Withhold timeout
* [ ] External components
    * [ ] Redis - cache inventory data
    * [ ] MQ - event handling
* [ ] Biz: query and check inventory available
* [ ] System optimistic
    * Connection pool tuning
    * PrepareStatement

## Optimize Solutions

* [ ] Scene: General row
    * Solution: simple update to make sure update success
* [ ] Scene: Hotspot row
    * [ ] Solution: advisory lock for less inventory
    * [ ] Solution: split single row into multiple for large inventory
        * Combination of display inventory id and real inventory id, which managed in item service
    * [ ] Solution: migrate hotspot row to standalone database
    * [ ] Solution: merge multiple deduct requests into single request
    * [ ] Solution: ratelimiter
* [ ] Scene: Multi-dimension inventory
