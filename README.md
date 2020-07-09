# File Hosting Site
[![CircleCI](https://circleci.com/gh/vpoletaev11/fileHostingSite.svg?style=svg)](https://circleci.com/gh/vpoletaev11/fileHostingSite)
[![Coverage Status](https://coveralls.io/repos/github/vpoletaev11/fileHostingSite/badge.svg?branch=master)](https://coveralls.io/github/vpoletaev11/fileHostingSite?branch=master)
### File hosting site intended for uploading/downloading files.<br>
Uploaded files will have: filename, file size, description, owner, category, upload date, rating.<br>
Site required that users should be registered.


# Project setup via docker
```shell
$ git clone https://github.com/vpoletaev11/fileHostingSite
$ cd fileHostingSite
$ docker-compose up
```
# Default project setup
## Step 1: Clone project

```shell
$ git clone https://github.com/vpoletaev11/fileHostingSite
```

## Step 2: Install dependencies
This project uses modules, because of this to install dependencies you just need to run tests
```shell
$ cd fileHostingSite
$ go test ./...
```

## Step 3: Configure database

```shell
$ mysql -u YOUR_MYSQL_USER < init.sql
$ YOUR_PREFERRED_REDACTOR main.go
Find this code: root:@tcp(mysql:3306) and reconfigure it. 
Syntax: username:password@connection_settings
```

## Step 4: Run project

```shell
$ go run main.go
```

## Step 5: Build project

```shell
$ go build main.go
```