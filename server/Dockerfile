FROM golang:1.20.2-alpine3.17 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /server

FROM alpine:3.17
WORKDIR /
COPY --from=build /server /server
ENTRYPOINT ["/server"]