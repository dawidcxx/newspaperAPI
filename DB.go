package main

import "github.com/jmoiron/sqlx"
import _ "github.com/lib/pq" // pg driver

//DB the database connection
//remember to initliaze it first. 
var DB *sqlx.DB
var err error

//InitDB initializes the database connection, panics if fails. 
func InitDB(con string) {
	DB, err = sqlx.Connect("postgres", con)
  if err != nil {
    panic(err.Error())
  }
}