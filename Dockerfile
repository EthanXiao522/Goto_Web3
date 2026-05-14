FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /server /server
COPY frontend/templates/ /templates/
COPY frontend/static/ /static/
EXPOSE 8080
CMD ["/server"]
