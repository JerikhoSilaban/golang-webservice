package main

import (
	"DTSGolang/Kelas2/Assignment7/controllers"
	"DTSGolang/Kelas2/Assignment7/routers"
)

var PORT = "127.0.0.1:8000"

func main() {
	controllers.StartDB()
	routers.StartServer().Run(PORT)
}
