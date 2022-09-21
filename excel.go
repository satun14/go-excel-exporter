package main

import (
	"database/sql"
	"github.com/xuri/excelize/v2"
	"strings"
	"time"
)

func makeFile(request *Request, rows *sql.Rows) error {
	f := excelize.NewFile()
	columns, _ := rows.Columns()

	fileName := strings.Join([]string{"export/export", time.Now().Format("_2006_01_02-15_04_05"), ".xlsx"}, "")

	sheet := "Sheet1"
	irow := 1
	//icol := 1
	cor := ""

	if err := drawFilters(f, sheet, &irow, request.Filters); err != nil {
		return err
	}

	if err := drawHeader(f, sheet, &irow, request.Fields); err != nil {
		return err
	}

	for rows.Next() {
		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(Scanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			return err
		}

		for icol, field := range request.Fields {
			for idx, column := range columns {
				if field.Name == column {
					var scanner = row[idx].(*Scanner)

					cor, _ = excelize.CoordinatesToCellName(icol+1, irow)

					err = f.SetCellValue(sheet, cor, scanner.Value)
					if err = rows.Err(); err != nil {
						return err
					}
				}
			}
		}

		irow++
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := setStyles(f, request, sheet, irow); err != nil {
		return err
	}

	if err := f.SaveAs(fileName); err != nil {
		return err
	}

	return nil
}

func drawFilters(f *excelize.File, sheet string, row *int, filters []Filter) error {
	for _, filter := range filters {
		cor, _ := excelize.CoordinatesToCellName(1, *row)

		err := f.SetCellValue(sheet, cor, filter.Label)
		if err != nil {
			return err
		}

		cor, _ = excelize.CoordinatesToCellName(2, *row)

		err = f.SetCellValue(sheet, cor, filter.Value)
		if err != nil {
			return err
		}

		*row++
	}

	*row++

	return nil
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

func setStyles(f *excelize.File, r *Request, sheet string, irow int) error {
	//BODY
	style, err := f.NewStyle(`{"font":{"family":"Calibri","size":11}}`)
	if err != nil {
		return err
	}

	corStart, _ := excelize.CoordinatesToCellName(1, 1)
	corEnd, _ := excelize.CoordinatesToCellName(len(r.Fields), irow)

	err = f.SetCellStyle(sheet, corStart, corEnd, style)
	if err != nil {
		return err
	}
	//

	//HEAD
	style, err = f.NewStyle(`{"border":[
	{"type":"left","color":"000000","style":1},
	{"type":"top","color":"000000","style":1},
	{"type":"bottom","color":"000000","style":1},
	{"type":"right","color":"000000","style":1}
	],"font":{"bold":true,"family":"Calibri","size":11}}`)
	if err != nil {
		return err
	}

	corStart, _ = excelize.CoordinatesToCellName(1, len(r.Filters)+2)
	corEnd, _ = excelize.CoordinatesToCellName(len(r.Fields), len(r.Filters)+2)

	err = f.SetCellStyle(sheet, corStart, corEnd, style)
	if err != nil {
		return err
	}
	//

	return nil
}
