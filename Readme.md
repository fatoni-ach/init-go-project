## Init DB
- buat database postgres baru
- restore file `backup_db.sql`

## Run Project

- download package
```
go mod download
```
- buat file config
```
cp copy.yml.example config.yml
```
- Jalankan aplikasi
```
 go run cmd/server/main.go
```

## Key Directory

* `cmd/server`: Main Golang
* `helpers`: All function Helper
* `pkg/delivery`: All Delivery Layer (This layer will act as the presenter. Decide how the data will presented. Could be as REST API, or HTML File, or gRPC whatever the delivery type.
This layer also will accept the input from user. Sanitize the input and sent it to Usecase layer.)
* `pkg/models`: All Models layer (Same as Entities, will used in all layer. This layer, will store any Object’s Struct and its method)
* `pkg/repository`: All Repository Layer (Repository will store any Database handler. Querying, or Creating/ Inserting into any database will stored here. This layer will act for CRUD to database only. No business process happen here. Only plain function to Database.)
* `pkg/usecase`: All Usecase Layer (This layer will act as the business process handler. Any process will handled here. This layer will decide, which repository layer will use. And have responsibility to provide data to serve into delivery. Process the data doing calculation or anything will done here.)
* `pkg/transformers`: All Struck to transforms data
* `services`: Service for cache, db, and log

## Clean Arc diagram:

![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)