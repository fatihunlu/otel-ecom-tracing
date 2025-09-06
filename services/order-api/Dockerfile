# ---- build
FROM golang:1.25 AS build
WORKDIR /src

# mod cache
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/order ./cmd/server

# ---- runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/order /app/order
EXPOSE 5002
ENTRYPOINT ["/app/order"]
