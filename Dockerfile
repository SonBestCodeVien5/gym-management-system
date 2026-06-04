ARG GO_VERSION=1.26.3

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/seed ./cmd/seed

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=build /out/server /app/server
COPY --from=build /out/seed /app/seed

USER app
EXPOSE 8080

CMD ["/app/server"]
