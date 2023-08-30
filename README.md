# Сервис динамического сегментирования пользователей

Микросервис для добавления пользователей в фича-эксперименты

# Используемые технологии:
- PostgreSQL
- Docker 
- Swagger
- Echo
- pgx

# Инструкция по запуску

Запустить сервис вместе с базой данных можно командой `make up` или `make up-attached`, если необходимо запустить докер в attached режиме

Этими командами запускается контейнеры с приложением и базой данных

# Swagger

Swagger доступен в браузере при запущенном контейнере приложения
`http://localhost/swagger/index.html` или в качестве json/yaml файла в директории `/docs`


# Примеры запросов и ответов


# Возникшие вопросы и их решения
**Откуда брать id пользователей?**

Решил, что id пользователей будут передаваться в запросе на создание пользователя, что логично учитывая то, что данное приложение -- микросервис.

**Что делать если при добавлении нескольких сегментов один из них не существует или неправильно указано expiration**

Решил, что в таком случае выполнение запроса будет фэйлится, так как возможно важно чтобы у пользователя был именно заданный набор фич и никакой другой

**Куда сохранять csv файлы, чтобы они были доступны по ссылке?**

Изначально хотел сохранять файлы в GoogleDrive, однако не успел разобраться с Google Cloud API, так что на данный момент файлы сохраняются в файловую систему контейнера
