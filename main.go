package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strings"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Excel Exporter Status - OK")
	})

	router.POST("/excel-export", func(c *gin.Context) {
		excel_export(c)
	})

	router.Run(":8090") // listen and serve on 0.0.0.0:8080
}

func excel_export(c *gin.Context) {
	var request Request

	if c.ShouldBind(&request) == nil {
		db, err := sql.Open("mysql",
			strings.Join([]string{
				request.Db.User,
				":",
				request.Db.Password,
				"@tcp(",
				request.Db.Host,
				":",
				request.Db.Port,
				")/",
				request.Db.Name,
				"?parseTime=true"}, ""))

		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		rows, err := db.Query(request.Sql)
		if err != nil {
			panic(err.Error())
		}

		columns, _ := rows.Columns()

		f := excelize.NewFile()

		fileName := strings.Join([]string{"export/export", time.Now().Format("_2006_01_02-15_04_05"), ".xlsx"}, "")

		sheet := "Sheet1"
		irow := 1
		//icol := 1
		cor := ""

		if err := drawHeader(f, sheet, &irow, request.Fields); err != nil {
			panic(err)
			//return c.JSON(http.StatusBadRequest, err)
		}

		for rows.Next() {
			row := make([]interface{}, len(columns))
			for idx := range columns {
				row[idx] = new(Scanner)
			}

			err := rows.Scan(row...)
			if err != nil {
				panic(err)
			}

			for icol, field := range request.Fields {
				for idx, column := range columns {
					if field.Name == column {
						var scanner = row[idx].(*Scanner)

						cor, _ = excelize.CoordinatesToCellName(icol+1, irow)

						err = f.SetCellValue(sheet, cor, scanner.Value)
						if err = rows.Err(); err != nil {
							panic(err)
							//return c.JSON(http.StatusBadRequest, err)
						}
					}
				}
			}

			irow++
		}

		if err = rows.Err(); err != nil {
			panic(err)
			//return c.JSON(http.StatusBadRequest, err)
		}

		if err := f.SaveAs(fileName); err != nil {
			panic(err)
			//return c.JSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, "Excel Done")
	}

	c.JSON(http.StatusBadRequest, "Dad request body")
}

func drawHeader(f *excelize.File, sheet string, row *int, fields []Field) error {
	for icol, field := range fields {
		cor, _ := excelize.CoordinatesToCellName(icol+1, *row)

		err := f.SetCellValue(sheet, cor, field.Label)
		if err != nil {
			return err
		}
	}

	*row++

	return nil
}
