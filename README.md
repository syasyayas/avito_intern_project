# Сервис динамического сегментирования пользователей

Микросервис для добавления пользователей в фича-эксперименты

# Используемые технологии:
- PostgreSQL
- Docker 
- Swagger
- Echo
- pgx

# Инструкция по запуску
Для запуска приложения понадобиться установленный docker-compose

Запустить сервис вместе с базой данных можно командой `make up` или `make up-attached`, если необходимо запустить докер в attached режиме

Этими командами запускаются контейнеры с приложением и базой данных

# Swagger

Swagger доступен в браузере при запущенном контейнере приложения 
`http://localhost/swagger/index.html` или в качестве json/yaml файла в директории `/docs`


# Примеры запросов и ответов

**Healthcheck**

`curl --location 'localhost:80/healthcheck'`


**Создание пользователя**

Запрос:
```
curl -i --location 'localhost:80/v1/user' \
--header 'Content-Type: application/json' \
--data '{"id": "123456"}'
```

Ответ:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```


**Получение пользователя и его сегментов**

Request:
```
curl --location --request GET 'localhost:80/v1/user' \
--header 'Content-Type: application/json' \
--data '{
    "id": "333"
}'
```

Response:

```
{
    "id": "333",
    "features": [
        {
            "slug": "AVITO_TEST_PERCENT100",
            "expires_at": "0001-01-01T00:00:00Z" // feature without expiration
        },
        {
            "slug": "AVITO_TEST1",
            "expires_at": "0001-01-01T00:00:00Z"
        },
        {
            "slug": "AVITO_TEST2",
            "expires_at": "0001-01-01T00:00:00Z"
        },
        {
            "slug": "AVITO_TEST3",
            "expires_at": "0001-01-01T00:00:00Z"
        },
        {
            "slug": "AVITO_TEST4",
            "expires_at": "2024-08-29T23:01:00Z" // feature with expiration
        }
    ]
}
```

**Удаление пользователя**

Request:
```
curl --location --request DELETE 'localhost:80/v1/user' \
--header 'Content-Type: application/json' \
--data '{
    "id":"333"
}'
```

Response :

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```

**Создание сегмента**

Request:

```
curl --location 'localhost:80/v1/feature' \
--header 'Content-Type: application/json' \
--data '{
    "slug":"AVITO_TEST4"
}'
```

Response:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```

**Создание сегмента с процентом пользователей**

Request:

```
curl --location 'localhost:80/v1/feature' \
--header 'Content-Type: application/json' \
--data '{
"slug":"AVITO_TEST_PERCENT_50"
}'
```


Response:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```


**Удаление сегмента**

Request:

```
curl --location --request DELETE 'localhost:80/v1/feature' \
--header 'Content-Type: application/json' \
--data '{
    "slug":"AVITO_TEST_PERCENT_50"
}'
```

Response:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```

**Добавление существующих сегментов пользователю**

Request:

```
curl --location 'localhost:80/v1/feature/features' \
--header 'Content-Type: application/json' \
--data '{
    "id":"333",
    "features":[
        {
            "slug":"AVITO_TEST_EXPIRE",
            "expires_at":"2024-08-29T23:01:00Z"
        },
        {
            "slug":"AVITO_TEST_NOEXPIRE"
        }
    ]
}'
```

Response:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```

**Удаление сегментов пользователя**

Request:

```
curl --location --request DELETE 'localhost:80/v1/feature/features' \
--header 'Content-Type: application/json' \
--data '{
    "id": "111",
    "features": [
        {
            "slug": "AVITO_TEST1"
        },
        {
            "slug": "AVITO_TEST2"
        }
    ]
}'
```
Response:

```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 18:08:11 GMT
Content-Length: 0
```

**Получение ссылки на csv файл с историей добавления/удаления сегментов у пользователя**

Request:

```
curl --location --request GET 'localhost:80/v1/history/export' \
--header 'Content-Type: application/json' \
--data '{
    "after":"2023-08-29T01:01:00Z",
    "before":"2023-08-29T23:01:00Z"
}'
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Thu, 31 Aug 2023 18:27:13 GMT
Content-Length: 78

{
    "url":"http://localhost:80/static/b0c98583-bded-4c09-9a18-d64ddebcdb8d.csv"
}

```

Пример файла лежит в папке /docs/example/

# Возникшие вопросы и их решения
**Откуда брать id пользователей?**

Решил, что id пользователей будут передаваться в запросе на создание пользователя, что логично, учитывая то, что данное приложение -- микросервис

**Что делать если при добавлении нескольких сегментов один из них не существует или неправильно указано expiration**

Решил, что в таком случае выполнение запроса будет фэйлится, так как, возможно, важно, чтобы у пользователя был именно заданный набор фич и никакой другой

**Куда сохранять csv файлы, чтобы они были доступны по ссылке?**

Изначально, хотел сохранять файлы в GoogleDrive, однако не успел разобраться с Google Cloud API, так что на данный момент файлы сохраняются в файловую систему контейнера с приложением
