FROM golang:alpine3.11 AS build
COPY . /app
WORKDIR /app
ENV GOPROXY="https://goproxy.cn"
RUN go build .


FROM alpine:3.11 AS prod
COPY --from=build /app/tgbot /app/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT [ "/app/tgbot" ]
