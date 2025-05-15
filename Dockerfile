FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/app

RUN chmod +x /backend

CMD ["/backend"]