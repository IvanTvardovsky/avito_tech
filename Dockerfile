FROM golang:1.20 AS build

WORKDIR /avito_tech
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /avito_tech/app

FROM alpine
EXPOSE 8080

COPY --from=build /avito_tech/app /app
COPY --from=build /avito_tech/config.yml /config.yml

ENTRYPOINT ["/app"]
