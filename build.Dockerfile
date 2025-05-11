# syntax=docker/dockerfile:1
FROM golang:1.24.3 AS build-stage

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/delve

FROM scratch
COPY --from=build-stage /app/delve /app/delve

CMD ["/app/delve"]


