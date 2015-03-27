#!/bin/bash

# Bind API HTTP server to this address
export IDP_ADDR="127.0.0.1:8000"
# Session's TTL (expires=now+TTL)
export IDP_SESSION_TTL=60
# Secret salt for hashing passwords (change before initial install, don't change afterwards)
export IDP_SECRET_SALT="842d7e1244b98f667f271a4e4d289772"

#
# MySQL
#
export IDP_DB_Driver="mysql"
export IDP_DB_DSN="root:@tcp(localhost:3306)/idp_dev?parseTime=true"

# Uncomment below for PostgreSQL
#export IDP_DB_Driver="postgres"
#export IDP_DB_DSN="postgres://alex:@localhost/idp_dev?sslmode=disable"

# Uncomment below for SQLite3
#export IDP_DB_Driver="sqlite3"
#export IDP_DB_DSN="/Users/alex/src/github.com/oleksandr/idp/db.sqlite3"

