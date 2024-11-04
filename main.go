package main

import (
	"example.com/RMS/db"
	"example.com/RMS/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	err:= db.InitDB()
	if err!=nil{
		panic("Database Error" + err.Error())
	}
	server := gin.Default()
	routes.Routes(server)
	server.Run("0.0.0.0:8080")

}
