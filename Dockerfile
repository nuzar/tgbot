FROM golang:1.16-alpine3.14 AS build
COPY . /app
WORKDIR /app
ENV GOPROXY="https://goproxy.cn"
RUN go build .


FROM alpine:3.14 AS prod
COPY --from=build /app/tgbot /app/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT [ "/app/tgbot" ]
