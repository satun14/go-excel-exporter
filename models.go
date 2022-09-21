package main

import "time"

type Request struct {
	Sql     string `json:"sql"`
	Db      Db     `json:"db"`
	Fields  []Field
	Filters []Filter
}

type Db struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type Filter struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Scanner struct {
	valid bool
	Value interface{}
}

func (scanner *Scanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

func (scanner *Scanner) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		if value, ok := src.(int64); ok {
			scanner.Value = value
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.Value = value
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.Value = value
			scanner.valid = true
		}
	case string:
		value := scanner.getBytes(src)
		scanner.Value = string(value)
		scanner.valid = true
	case []byte:
		value := string(scanner.getBytes(src))
		scanner.Value = value
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.Value = value.Format("02.01.2006")
			scanner.valid = true
		}
	case nil:
		scanner.Value = nil
		scanner.valid = true
	}
	return nil
}
