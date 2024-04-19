## Распределенный вычислитель арифметических выражений
### [Ссылка на telegram автора](https://t.me/iyozy)

## Описание
Комплекс сервисов илюстрирующих распределенную модель вычислений и основным паттерны, которые могут быть использованы
для ее реализации.

## В основе лежат следующиет требования
Пользователь хочет считать арифметические выражения.
Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго.
Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна.
Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности.
Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин".
Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"?
Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

## Разворачивание
Для запуска требуется выполнить сначала команду
```
docker compose up -d
```
Потом создать структуру таблиц с помощью скрипта миграции
Потребуется инструмент migrate [source](https://github.com/golang-migrate/migrate)
```
migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" up
```
Для удаления таблиц
```
migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" down
```

## Особенности реализации
В docker compose описано создание системы из 6 контейнеров. 
1) calculator-api. http сервер, принимает запросы от клиента с заданиями и от агентов вычислителей с посчитанными заданиями и heartbeat.
Поддерживает аутентификацию с помощью jwt-токенов. Под авторизацию сделана часть работы, пользователи хранятся с дополнительным признаком Admin.
В будещем эта часть будет доработана.
2) message-broker. Брокер сообщений RabbitMQ, используется для передачи сообщений агентам-вычислителям.
3) calculator_db. База данных на Postgres. В нее вносится информация о посчитанных выражениях, подвыражениях, которые находятся в процессе рассчета,
пользовательские настройки длительностей операций и региструются heartbeat агентов вычислителей.
4) calculator-dbagent. Сервис, созданный для работы с базой данных, агенты и api общаются с ним по grpc, сам сервис работает с базой напрямую.
5) calculator-agentN. Агент вычислитель, достает подвыражение из очереди и передает рабочему внутри сервиса на вычисление, количество одновременно допустимых
рабочих регулируется переменной среды. В данном примере поднимаются 2 агента.

## Использование
### Видео -> [Демонстрация сценария использования](https://www.loom.com/share/ec6a52466bf84f6da9ca72947d2fbe74?sid=0096f914-f553-4066-888b-7ac36adbb974)

Сервис api принимает следующие запросы.

### Группа Auth
1. Регистрация нового пользователя
```
curl --location --request POST 'http://localhost:8080/api/v1/register' --header 'Content-Type: application/json' --data '{
    "login": "user22",
    "password": "1+2+3+4"
}'
```
Регистрируем нового пользователя, проверка на уникальность логина осуществляется на уровне базы, в будущем лучше сделать ошибку более приятной для пользователя.
Пароль в базе храним в зашифрованном виде, чтобы получив доступ к базе нельзя было утащить пароль.
Пример ответа, содержит id созданного пользователя
```
{"id":2,"status":200}
```
2. Login.
```
curl --location --request POST 'http://localhost:8080/api/v1/login' --header 'Content-Type: application/json' --data '{
    "login": "user22",
    "password": "1+2+3+4"
}'
```
Проводим проверку логина пароля, если все успешно, генерируем для него access/refresh токены. Время жизни прописывается в конфиге.
Пример ответа, refresh прописываем в cookie
```
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTQwNzIsImlhdCI6MTcxMzU1MDQ3MiwidXNlcklkIjo0fQ.gyVEYqy5VxG9MRInKv00K-FkGpyKrFaU6xVhaICAXp4"
}
```
3. Refresh. 
Запрос, которым мы обновляем пару токенов, если время жизни access токена вышло. Когда выйдет время жизни refresh token прийдется логинится заново.
```
curl --location --request POST 'http://localhost:8080/api/v1/refresh' --header 'Cookie: refresh-token=128cb3b45f1ddf29d9a7b7c208657825440e3363bf00d1569385b3e2196b882c' --data ''
```
Пример ответа, refresh прописываем в cookie
```
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTQwNzIsImlhdCI6MTcxMzU1MDQ3MiwidXNlcklkIjo0fQ.gyVEYqy5VxG9MRInKv00K-FkGpyKrFaU6xVhaICAXp4"
}
```
### Группа Expression
1. Запрос на постановку выражения на вычисление.
```
curl --location --request POST 'http://localhost:8080/expression' --header 'Content-Type: application/json' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTgyMjYsImlhdCI6MTcxMzU1NDYyNiwidXNlcklkIjoxfQ.5tWA43_qYK1catJjEq_0WwJeu9PCs3H8hOyJX2NWVKQ' --header 'Cookie: refresh-token=128cb3b45f1ddf29d9a7b7c208657825440e3363bf00d1569385b3e2196b882c' --data '{
    "req_uid": "23332e",
    "expression": "1+2+3+4"
}'
```
Допустимы выражения, хотя бы с одним знаком операции, не более 50 символов в длину(регулируется параметром в конфиге), состоящие из
чисел(целые и с плавающей точкой) и знаков операций `+-/*`. Скобки, двойные знаки операций(-*, +-, -+ и так далее) недопустимы, планирую добавить позднее.
req_uid - Это уникальный идентификатор запроса, нужен для того, чтобы если пользователь не получил ответ содержащий созданное для расчета выражение, он мог
повторным запросом его найти. Если под данным uid уже существует выражение для данного пользователя, то сервис просто вернет его пользователю, между пользователями req_uid запросов могут пересекаться.

Пример ответа
```
{
    "expression": {
        "expression_id": 2,
        "user_id": 1,
        "req_uid": "23332e",
        "expression_value": "1+2+3+4",
        "result": 10,
        "status": "Ok",
        "error_msg": "",
        "created_at": "2024-04-19T19:27:03Z",
        "eval_started_at": "2024-04-19T19:27:03Z",
        "eval_finished_at": "2024-04-19T19:27:04Z"
    },
    "status": 200
}
```

2. Запрос на отображение всех выражений. Подфильтровывается userId из токена
```
curl --location --request GET 'http://localhost:8080/expression' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTgyMjYsImlhdCI6MTcxMzU1NDYyNiwidXNlcklkIjoxfQ.5tWA43_qYK1catJjEq_0WwJeu9PCs3H8hOyJX2NWVKQ' --header 'Cookie: refresh-token=128cb3b45f1ddf29d9a7b7c208657825440e3363bf00d1569385b3e2196b882c' --data ''
```
Выражения могут быть в процессе расчета, рассчитаны с ошибкой(деление на 0), расчитаны без ошибки, либо еще не взяты агентами.

### Группа Duration
1. Запрос на добавление замедления для той или иной математической операции. Так как мы имитируем нагрузку на вычислителях. Каждую из четырех допустимых
математических операций можно замедлить. Замедление передается в ms. Допустимы значения от 0 до 65535
Если длительность уже зарегистрированна, ее значение будет перезаписано на новое переданное в запросе. По умолчанию замедление 0.
Все пользователи могут менять эти парамерты и они общие для всех запросов, в будущем лучше перевести на пользователя с админскими правами.
```
curl --location --request POST 'http://localhost:8080/duration' --header 'Content-Type: application/json' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTgyMjYsImlhdCI6MTcxMzU1NDYyNiwidXNlcklkIjoxfQ.5tWA43_qYK1catJjEq_0WwJeu9PCs3H8hOyJX2NWVKQ' --header 'Cookie: refresh-token=128cb3b45f1ddf29d9a7b7c208657825440e3363bf00d1569385b3e2196b882c' --data '{
    "+":1000,
    "-":1000,
    "*":1000,
    "/":10000
}'
```
2. Запрос на отображение текущих длительностей по операциям.
```
curl --location --request GET 'http://localhost:8080/duration' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTgyMjYsImlhdCI6MTcxMzU1NDYyNiwidXNlcklkIjoxfQ.5tWA43_qYK1catJjEq_0WwJeu9PCs3H8hOyJX2NWVKQ' --data ''
```
Пример ответа:
```
{
    "operation_durations": [
        {
            "operation_name": "+",
            "operation_duration_ms": 1000,
            "created_at": "2024-04-19T19:35:17Z",
            "updated_at": "2024-04-19T19:35:33Z"
        },
        {
            "operation_name": "-",
            "operation_duration_ms": 1000,
            "created_at": "2024-04-19T19:35:17Z",
            "updated_at": "2024-04-19T19:35:33Z"
        },
        {
            "operation_name": "*",
            "operation_duration_ms": 1000,
            "created_at": "2024-04-19T19:35:17Z",
            "updated_at": "2024-04-19T19:35:33Z"
        },
        {
            "operation_name": "/",
            "operation_duration_ms": 10000,
            "created_at": "2024-04-19T19:35:17Z",
            "updated_at": "2024-04-19T19:35:33Z"
        }
    ],
    "status": 200
}
```
### Группа Agents
1. Запрос на возвращение всех Агентов. Возращает список зарегистрированных агентов-вычислителей, с их именем, датой создания и heartbeat.
```
curl --location --request GET 'http://localhost:8080/agent' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NTgyMjYsImlhdCI6MTcxMzU1NDYyNiwidXNlcklkIjoxfQ.5tWA43_qYK1catJjEq_0WwJeu9PCs3H8hOyJX2NWVKQ' --header 'Cookie: refresh-token=128cb3b45f1ddf29d9a7b7c208657825440e3363bf00d1569385b3e2196b882c'
```

Пример ответа
```
{
    "agents": [
        {
            "agent_unique_name": "AGENT_1",
            "created_at": "2024-04-19T18:43:33Z",
            "last_heartbeat_at": "2024-04-19T19:40:49Z"
        },
        {
            "agent_unique_name": "AGENT_2",
            "created_at": "2024-04-19T18:43:32Z",
            "last_heartbeat_at": "2024-04-19T19:40:49Z"
        }
    ],
    "status": 200
}
```

## Переменные окружения
### calculator-api
Соль и секрет для пароля пользователя и JWT токенов.
```
  AUTH_SALT: KKK
  AUTH_SECRET: XXX
```
### dbagent
Параметры для подключения к базе данных
```
   DB_HOST=postgres
   DB_PORT=5432
   DB_USERNAME=postgres
   DB_NAME=calculator_db
   DB_SSLMODE=disable
   DB_PASSWORD=postgres
```

### agent
Имя агента и количество парралельных вычислений, которые он может делать
```
  AGENT_NAME=AGENT_1
  AGENT_MAX_PROCESSES=5
```
### calculator_db
```
  POSTGRES_PASSWORD=postgres
  POSTGRES_USER=postgres
  POSTGRES_DB=calculator_db
```

### Особенности организации очереди
В брокере мы создаем 2 очереди, одна для работы, другая для сообщений, которые не получили подтверждения обработки после вычитывания, эта очередь должна замедлять повторное вычитывание на какое-то время,
чтобы агенты не сильно спамили лог, если что-то случилось с dbagent или базой, по-хорошему это время должно быть таким, чтобы хватило времени разобраться с аварией, и ничего не обрабатывалось в этот момент.
Подразумевается следующий сценарий, хотя для данного проекта он выглядит избыточным, так как данная исключительная ситуация покрывается другим сценарием
самовосттановления. Агент вычитал задачу из очереди, обработал ее и пытаясь сохранить ответ получил ошибку. В этом случае он ставит сообщению в очереди Nack,
что приводит к тому, что сообщение попадает в очередь исключений на 20 секунд(параметр задается в конфиге), а затем возвращается в основную очередь. Где может
быть взят любым свободным вычислителем, вычисление будет проведено заново и будет завново предпринята попытка сообщить результат dbagent. 

### Схема взаимодействия сервисов
![Схема взаимодействия сервисов](schema.png "Схема взаимодействия сервисов")

### Описание основных сценариев
1. С помощью POST запроса создаем пользователя. Далее пользователь логинится полчает пару токенов. С помощью токена в дальнейшем подписывает все свои запросы.
Клиент передает выражение. Api разбивает выражение на подвыражения отправляет их dbagent, который сохраняет в БД. Отправляет подвыражения, которые могут быть посчитаны
независимо в очередь, оттуда они вычитываются агентами, агент передает dbagent агенту информацию, что он взял выражение, а чуть позже результат расчетов. Dbagent в свою очередь сохраняет этов все в БД.
Дополнительно агент посылает Heartbeatы, которые тоже сохраняет dbagent. Если было  подсчитано последнее повыражение выражения, dbagent акуумулирует информацию в основном выражении и меняет ему статус на посчитанное.
2. Агент потерял связь с dbagent, если при этом у api связь осталась, то спустя какое-то время опеределенное в конфиге, api поймет, что последний hb от агента был давно и поскидывает его пометку со всех задач,
они опять попадут в очередь и могут быть взяты другими агентами.
3. Агент упал в момент рассчета, подвыражение было вычитано из очереди, в этом случае подвыражение не получит подтвержения, отправится в очередь ожидания и спустя время, определенное в конфиге, появится в основной очереди.
Это защита от спама ошибками.
4. Упала база, выражения перестанут создавать, регистрироваться, будут попытки разбора из очереди, но учитывая, что нельзя записать ответ, агенты не будут подтвержать вычитывания и выражения будут возращаться в очередь.
5. Умер DBAgent, ситуация аналогична ситуации с базой, возможно нужно поднимать для защиты от такой ситуации их несколько.
6. Умерло API, агенты смогут доразобрать, досчитать выржаения, которые были в очереди, потом останутся без работы.


### Что еще хотелось бы сделать
1. Написать нормальные тесты, не хватило времени, смог покрыть не так много.
2. Встроить migrate в образы. Добавить замедлители для старта зависимых сервисов.
3. Подержать механизм авторизации, через создание админских ролей. Единственное пока не могу понять, кто заведет админа, возможно скриптами при миграции.
4. Облагородить валидацию и сделать красивые ошибки для пользователей, есть ряд ошибок, которые прилетают с базы и хотя они перехватываются, текст выглядит некрасиво. 
5. Удалять из таблицы с подвыражениями, подвыражения, для посчитанных выражений. Они уже ни к чему, только пухнущая таблица, скорее как логи на первом этапе.

