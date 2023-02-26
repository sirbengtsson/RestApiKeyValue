package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"strconv"
)

type KeyValueStore interface {
	setValue(key string, value string) error
	readValue(key string) (string, error)
	deleteValue(key string) error
	getAllKeys() ([]string, error)
}

type SaveValueRequest struct {
	Value string `json:"value"`
}

func main() {
	persistent := os.Getenv("IS_PERSISTENT") == "true"
	var keyValueStore KeyValueStore

	if persistent {
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			panic(err)
		}

		keyValueStore = &PersistentStore{
			redisClient: redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_HOST"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       redisDb,
			}),
		}
	} else {
		keyValueStore = &MemoryStore{}
	}

	router := gin.Default()

	router.POST("/api/v1/cache/:key", func(c *gin.Context) {
		var requestBody = SaveValueRequest{}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Cannot parse request body",
			})
		}
		fmt.Println(requestBody)

		err := keyValueStore.setValue(c.Param("key"), requestBody.Value)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Cannot save key because %s", err.Error()),
			})
		}
	})

	router.GET("/api/v1/cache/:key", func(c *gin.Context) {
		val, err := keyValueStore.readValue(c.Param("key"))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"value": val})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Cannot get key because %s", err.Error()),
			})
		}
	})

	router.DELETE("/api/v1/cache/:key", func(c *gin.Context) {
		err := keyValueStore.deleteValue(c.Param("key"))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Cannot get key because %s", err.Error()),
			})
		}
	})

	router.GET("/api/v1/cache/keys", func(c *gin.Context) {
		keys, err := keyValueStore.getAllKeys()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"keys": keys})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Cannot get key because %s", err.Error()),
			})
		}
	})

	router.Run()
}
