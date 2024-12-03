# Portal-RMS

Service used for comapnies to manage the Job Openings and Applications (JOB PORTAL)

## Technologies

* Go - 1.23
* Gin
* GORM
* Postgresql

## Setup
```
git clone
cd portalrms
 ```

One time db setup

``Create Database rms with user rms and password rms``

Env File

``create .env file if not present and add below content``

```

DB_HOST = 127.0.0.1
DB_PORT = 5432
DB_USER =  rms
DB_PASSWORD =  rms
DB_NAME = rms
DB_SSLMODE = require
AUTH_KEY = "#123&456VR" used in JSON WEB TOKENS
```

## Build
 ``make build``

## Run locally

* Database migration is handled using GORM.

`go run /main.go` 
