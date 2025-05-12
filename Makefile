APP_NAME=kiraform-service
ENTRY_POINT=src/infras/entry/main.go

# define commands
run:
	go run ${ENTRY_POINT}

build:
	go run build -o ${APP_NAME} ${ENTRY_POINT}