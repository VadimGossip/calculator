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
Для удаления таблицу
```
migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" down
```

## Особенности реализации
В docker compose описано создание системы из 5 контейнеров. 
1) calculator-api. http сервер, принимает запросы от клиента с заданиями и от агентов вычислителей с посчитанными заданиями и heartbeat
2) message-broker. Брокер сообщений RabbitMQ, используется для передачи сообщений агентам-вычислителям.
3) calculator_db. База данных на Postgres. В нее вносится информация о посчитанных выражениях, подвыражениях, которые находятся в процессе рассчета,
пользовательские настройки длительностей операций и региструются heartbeat агентов вычислителей.
4) calculator-agentN. Агент вычислитель, достает подвыражение из очереди и передает 

## Использование



Сервис принимает POST запросы на https://{subproject_name}-{env_name}.etservice.net/update_acc

Формат запроса

```
type UpdateAccRequest struct {
	Id               uint64 `json:"id" binding:"required"`
	AccId            uint64 `json:"acc_id" binding:"required"`
	AccAllowOutgoing *int   `json:"acc_allow_outgoing,omitempty"`
	AccAllowIncoming *int   `json:"acc_allow_incoming,omitempty"`
}
```

Формат ответа

```
type CommonResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}
```

Валидация
Обязательнымы являются
id - id запроса
acc_id - account id параметры, которого мы будем менять.
Также должно быть обязательно заполнено одно из полей состояния, либо acc_allow_outgoing, либо acc_allow_incoming, допустима передача обоих полей.

### Health
Формат ответа
```
type HealthResponse struct {
	Status       string    `json:"status"`
	AppStartedAt time.Time `json:"app_started_at"`
}
```
В ответе мы получаем инормацию о том, что http сервер сервиса поднят/не поднят и дату старта сериса.

### Metrics
В качестве метрики мы логируем дополнительно ответы на все запросы.

Пример:
```
tj_acl_alaris_api_response{code="200",error="",method="POST",path="/update_acc"} 1
tj_acl_alaris_api_response{code="400",error="Parse request error: Key: 'UpdateAccRequest.AccId' Error:Field validation for 'AccId' failed on the 'required' tag",method="POST",path="/update_acc"} 1
```

## Переменные окружения
APP_HTTP_PORT: "8080"
ALARIS_API_BASEURL: "http://xxxx.com:port"
ALARIS_API_AUTHKEY: "key_value"


## Распределенный вычислитель арифметических выражений



### В основе лежит следующее тестовое задание



### Реализация

В моей реализации комплекс состоит из некоторого сервиса api, блокера сообщений RabbitMQ, базы данных Postgres, и сервиса вычислителя, и