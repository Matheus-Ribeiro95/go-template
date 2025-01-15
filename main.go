package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"matheusribeiro95/template/lib"
	"net/http"
	"os"
)

type API_POST_INSERT_BODY struct {
	Name string
}

type API_POST_DELETE_BODY struct {
	Id string
}

type API_POST_UPDATE_BODY struct {
	Id   string
	Name string
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	LIB.Upgrader()

	var SCYLLADB_URL = os.Getenv("SCYLLADB_URL")
	if SCYLLADB_URL == "" {
		SCYLLADB_URL = "localhost"
	}
	var cluster = gocql.NewCluster(fmt.Sprintf("%s:9042", SCYLLADB_URL))

	var session, clusterError = cluster.CreateSession()
	if clusterError != nil {
		panic("Failed to connect to cluster")
	}

	defer session.Close()

	router := gin.Default()

	router.Static("/static", "./html")

	router.GET("/api/getAll", func(c *gin.Context) {
		response := LIB.SelectAll(session, logger)

		if response == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    response,
			})
		}
	})

	router.POST("/api/insert", func(c *gin.Context) {
		var data API_POST_INSERT_BODY
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})

			return
		}

		response, execute := LIB.InsertQuery(data.Name, session, logger)

		if execute == false {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})
		} else {
			LIB.NewUser(response)

			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    response,
			})
		}
	})

	router.POST("/api/delete", func(c *gin.Context) {
		var data API_POST_DELETE_BODY
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})

			return
		}

		response, execute := LIB.DeleteQuery(data.Id, session, logger)

		if execute == false {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})
		} else {
			LIB.DeleteUser(response)

			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    response,
			})
		}
	})

	router.POST("/api/update", func(c *gin.Context) {
		var data API_POST_UPDATE_BODY
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})
			return
		}

		response, execute := LIB.UpdateQuery(data.Id, data.Name, session, logger)

		if execute == false {
			c.JSON(http.StatusOK, gin.H{
				"message": "erro",
				"data":    false,
			})
		} else {
			LIB.UpdateUser(response)

			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    response,
			})
		}
	})

	router.GET("/ws", func(c *gin.Context) {
		LIB.Echo(c.Writer, c.Request)
	})

	router.Run()
}
