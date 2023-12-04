package main

import ("github.com/gin-gonic/gin"

"github.com/milantracy/playground/internal/ws")

func main() {
	router := gin.Default()
	router.GET("/albums", ws.GetAlbums)

	router.Run("localhost:8080")
}
