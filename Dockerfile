FROM golang:1.21.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /main

EXPOSE 8080

# Run
CMD ["/main"]
