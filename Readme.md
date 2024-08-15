## Init DB
- buat database postgres baru
- restore file `init-project/backup_db.sql`
- run sql `init-project/seeder.sql`
- collection postman `init-project/Asset Management.postman_collection.json`

## Run Project

- download package
```
go mod download
```
- buat file config
```
cp .env.example .env
```
- Jalankan aplikasi
```
 go run cmd/server/main.go
```
## How to Deploy
- Set `.env` file based `.env.example`  file
```
cp .env.example .env
```
- Adjust the value of .`env` file
- Run this command to compile app
```
go build cmd/server/main.go
``` 
- Run the app
```
./main
```

## Running service using systemd service
- Make new systemd service
```
sudo nano /etc/systemd/system/jobs.service
```
- Add this code and then save
```
[Unit]
Description=gowebapi

[Service]
Type=simple
Restart=always
RestartSec=10s
WorkingDirectory=/var/www/html/api-deepface-recognition/(Working direcroty app)
ExecStart=/var/www/html/api-deepface-recognition/main(compiled file)

[Install]
WantedBy=multi-user.target
```
- Run the service
```
sudo systemctl start jobs
sudo systemctl enable jobs
```
- To stop the service run this command :
```
sudo systemctl stop jobs
```
- To remove the service
```
sudo systemctl disable jobs
```
- You can check the service in the web browser using port which you already set.

<br>

## Deploy using Dockerfile
- Build Image
```
docker build -t go-project-image .
```
- Run Image in Container
```
docker run -it --rm --name go-project-container go-project-image
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