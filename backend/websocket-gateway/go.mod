module github.com/todoist/backend/websocket-gateway

go 1.21

require (
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/todoist/backend/pkg v0.0.0
)

replace github.com/todoist/backend/pkg => ../pkg

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
