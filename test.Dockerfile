FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build

FROM adoptopenjdk/openjdk11:alpine-jre
WORKDIR /app
RUN apk add iptables
COPY . .
COPY --from="build" /app/craftignite .
EXPOSE 25565:25565
CMD ["./craftignite"]
