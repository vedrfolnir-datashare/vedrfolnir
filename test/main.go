package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	TestEncrypt()
	TestSearch()
	TestIndexConstruction()
	TestIndexVerifyTimeSpaceAll()
}
