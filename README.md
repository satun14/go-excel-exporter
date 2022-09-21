# Учебный проект по Golang - ExcelExporter

Проект для экспорта данных из БД в файл Excel, по заданному запросу и набору полей.

Веб-сервер получает на вход json запрос со следующими данными:
```
{
  "db": {
    "host": "localhost",
    "port": "3306",
    "name": "world",
    "user": "root",
    "password": "qwerty"
  },
  "sql": "select * from country where Region = 'Southern Europe'",
  "fields": [
    {
      "name": "Name",
      "type": "string",
      "label": "Название"
    },
    {
      "name": "Code",
      "type": "string",
      "label": "Код"
    }
  ],
  "filters": [
    {
      "label": "Регион",
      "value": "Southern Europe"
    }
  ]
}
```

Где все параметры, кроме блока __filters__ являются обязательными.

### Запуск
```
$ go run .
```

### Настройка
Можно указать хост и порт сервера через __.env__ файл.