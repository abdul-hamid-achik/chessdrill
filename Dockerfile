FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY go.mod go.sum ./
RUN go mod download
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN templ generate
RUN npm run build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/chessdrill ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/chessdrill /app/chessdrill
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["/app/chessdrill"]
