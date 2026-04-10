# stage 1: build
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cardinal .
# install goose binary
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

# stage 2 - final image
FROM alpine:latest 

RUN apk add --no-cache ca-certificates

WORKDIR /app 

# copy binary and goose
COPY --from=builder /app/cardinal .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# copy static files
COPY --from=builder /app/index.html .
COPY --from=builder /app/dashboard.html .
COPY --from=builder /app/profile.html .
COPY --from=builder /app/verified.html .
COPY --from=builder /app/js ./js
COPY --from=builder /app/css ./css

# copy migrations
COPY --from=builder /app/sql/schema ./sql/schema

# copy entrypoint
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh 

EXPOSE 8080 

ENTRYPOINT [ "./entrypoint.sh" ]