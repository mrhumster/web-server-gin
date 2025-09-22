package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album reprtesent data about a record album
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Selected", Artist: "Amelie Lens", Price: 15.99},
	{ID: "2", Title: "A Moment Apart", Artist: "ODESZA", Price: 12.99},
	{ID: "3", Title: "Exhale", Artist: "Charlotte de Witte", Price: 14.99},
	{ID: "4", Title: "The Age of Love", Artist: "Age of Love", Price: 10.99},
	{ID: "5", Title: "Raveolution", Artist: "Dax J", Price: 13.99},
	{ID: "6", Title: "In Silence", Artist: "Nina Kraviz", Price: 16.99},
	{ID: "7", Title: "Planet-X", Artist: "Jeff Mills", Price: 11.99},
	{ID: "8", Title: "The Black Madonna", Artist: "We Still Believe", Price: 14.49},
	{ID: "9", Title: "Hypercolour", Artist: "Shanti Celeste", Price: 12.49},
	{ID: "10", Title: "Aether", Artist: "Ben Klock", Price: 15.49},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "healthy")
}

func main() {
	router := gin.Default()
	router.GET("/health", getHealth)
	router.GET("/albums", getAlbums)
	router.Run(":8080")
}
