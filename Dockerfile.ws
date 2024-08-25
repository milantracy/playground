# syntax=docker/dockerfile:1

# update here if go version is updated for the project.
FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Move the necessary files only.
COPY api/ws/ ./api/ws/
COPY internal/ws/ ./internal/ws/
COPY cmd/ws/ ./cmd/ws/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/ws-albums ./cmd/ws/

CMD ["./bin/ws-albums"]


