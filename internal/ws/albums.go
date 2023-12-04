package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milantracy/playground/api/ws"
)

var albums = map[string]ws.Album{
	"1": {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	"2": {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	"3": {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// GetAlbums returns the list of all albums.
func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// GetAlbumByID returns a album whose ID matches the given ID parameter.
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	album, exist := albums[id]
	if !exist {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found!"})
	}
	c.IndentedJSON(http.StatusOK, album)
}

// AddAlbum adds an album from the request.
func AddAlbum(c *gin.Context) {
	var newAlbum ws.Album

	if err := c.BindJSON(&newAlbum); err == nil {
		albums[newAlbum.ID] = newAlbum
		c.IndentedJSON(http.StatusCreated, newAlbum)
	}
}
