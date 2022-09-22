package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

func main() {
	godotenv.Load()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Excel Exporter Run")
	})

	router.POST("/excel-export", func(c *gin.Context) {
		excelExport(c)
	})

	serverHost := os.Getenv("HOST")
	serverPort := os.Getenv("PORT")

	if serverHost != "" && serverPort != "" {
		router.Run(strings.Join([]string{serverHost, ":", serverPort}, ""))
	} else {
		router.Run()
	}
}

func excelExport(c *gin.Context) {
	var request Request

	if c.ShouldBind(&request) == nil {
		rows, err := execSqlQuery(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := makeFile(&request, rows); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, "Excel Done")
	} else {
		c.JSON(http.StatusBadRequest, "Bad request body")
	}
}
