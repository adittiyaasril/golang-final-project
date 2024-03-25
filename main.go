package main

import (
	"final-project/database"
	"final-project/routers"
)

func main() {
	database.StartDB()
	r := routers.StartApp()
	r.Run(":8080")
}
