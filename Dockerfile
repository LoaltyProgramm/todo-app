FROM golang:1.23.6

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd /app/cmd
COPY internal /app/internal
COPY web /app/web

RUN GOOS=linux GOARCH=amd64 go build -o todo-app /app/cmd/main.go

EXPOSE 7540

CMD [ "./todo-app" ]

