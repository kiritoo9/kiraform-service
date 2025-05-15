APP_NAME=kiraform-service
ENTRY_POINT=src/infras/entry/main.go

# define commands
swag:
	swag init -g ${ENTRY_POINT}

run:
	go run ${ENTRY_POINT}

build:
	go run build -o ${APP_NAME} ${ENTRY_POINT}