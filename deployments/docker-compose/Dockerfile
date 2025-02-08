FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api cmd/app/main.go

#---

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/api .
COPY --from=builder /app/cmd/app/.env .

RUN ls -lah

ENTRYPOINT [ "./api" ]