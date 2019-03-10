FROM golang:latest AS build-env

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod=vendor -o explorer


FROM alpine

WORKDIR /app
COPY --from=build-env /src/explorer .

EXPOSE 8080

CMD [ "./explorer" ]
