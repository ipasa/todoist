module github.com/todoist/backend/task-service

go 1.21

require (
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/todoist/backend/pkg v0.0.0
)

replace github.com/todoist/backend/pkg => ../pkg

require (
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
