# File Hosting Site

### File hosting site intended for uploading/downloading files.<br>
Uploaded files will have: filename, file size, description, owner, category, upload date, rating.<br>
Site required that users should be registered.


# Project setup
## Step 1: Clone project

```shell
$ git clone https://github.com/vpoletaev11/fileHostingSite
```

## Step 2: Configure database

```shell
$ mysql -u YOUR_MYSQL_USER < ~/fileHostingSite/database/databaseStructure/createDB.sql
$ cd fileHostingSite
$ YOUR_PREFERRED_REDACTOR main.go
Find this code: perdator:@tcp(localhost:3306) and reconfigure it. 
Syntax: username:password@connection_settings
```

## Step 3: Run project

```shell
$ go run main.go
```

## Step 4: Build project

```shell
$ go build.main.go
```

