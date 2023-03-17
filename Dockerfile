FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build -o account_api .

#Build
FROM ubuntu:latest
WORKDIR /
COPY --from=build /app/account_api /account_api
EXPOSE 9100
CMD ["/account_api"]