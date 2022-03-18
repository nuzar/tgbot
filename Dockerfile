FROM docker.io/library/golang:1.18-alpine3.15 AS build
COPY . /app
WORKDIR /app
RUN go build -buildvcs=false .


FROM docker.io/library/alpine:3.15 AS prod
COPY --from=build /app/tgbot /app/
WORKDIR /app
EXPOSE 8080
CMD [ "/app/tgbot" ]
