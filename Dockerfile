FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build -o account_api .

#Build
FROM ubuntu:latest
WORKDIR /app
COPY --from=build /app/.env /app
COPY --from=build /app/account_api /app/account_api

CMD ["/app/account_api"]