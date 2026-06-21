From golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server-app ./cmd

EXPOSE 8080

CMD ["./server-app"]