FROM golang:1.19-alpine as build

WORKDIR /go/src/app

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY . .

ENV APP_NAME=chatelegram

RUN ["go", "build", "-o", "build/${APP_NAME}"]

FROM build as dev

CMD ["go", "run", "."]

FROM alpine:3.14.1 as prod

RUN addgroup -S app && adduser -S -G app app
USER app

WORKDIR /home/app/

COPY --from=build /go/src/app/build/${APP_NAME} ./

CMD ["./${APP_NAME}"]
