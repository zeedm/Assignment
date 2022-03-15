# Assignment

This is an assignment about an order backend system that allows users to order products from vendors online that also tracks vendor inventory. This system is to cater for both users and vendors.

## Starting the project

To run the project, you must [clone the project first](https://www.atlassian.com/git/tutorials/setting-up-a-repository/git-clone)
You must also have [SQL Server](https://docs.microsoft.com/en-us/sql/database-engine/install-windows/install-sql-server?view=sql-server-ver15) in your local machine and [retore the backup file](https://docs.microsoft.com/en-us/sql/relational-databases/backup-restore/quickstart-backup-restore-database?view=sql-server-ver15#restore-a-backup)
Then, you can "go run ." in the Terminal to start the project.

```console
go run .
```

## Libraries

[go-mssqldb](github.com/denisenkom/go-mssqldb): to connect to SQL Server

[gorilla/mux](): to create APIs and handle requests

[golang-jwt/jwt](github.com/golang-jwt/jwt): to create jwt

[gorilla/sessions](github.com/gorilla/sessions): to handle sessions

[stretchr/testify](github.com/stretchr/testify): to create test

[go-sqlmock](github.com/DATA-DOG/go-sqlmock): to mock db for testing

## Architecture

The project uses basic MVC architecture pattern, Golang as the coding language, and SQL Server to store data.