# build stage
FROM golang:1.10 as build-env

WORKDIR /go/src/github.com/SmotrovaLilit/traefik-auth-botforbidden

RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . .

RUN cd /go/src/github.com/SmotrovaLilit/traefik-auth-botforbidden/cmd/auth && CGO_ENABLED=0 GOOS=linux go build -a -o ../../main

# final stage
FROM alpine:latest

ENV PORT=80

LABEL maintainer="Smotrova Lilit <lilit2990@gmail.com>"

WORKDIR /app
EXPOSE 80

COPY --from=build-env /go/src/github.com/SmotrovaLilit/traefik-auth-botforbidden/main .
COPY --from=build-env /go/src/github.com/SmotrovaLilit/traefik-auth-botforbidden/entrypoint.sh .

RUN chmod -R 755 /app/entrypoint.sh

ENTRYPOINT sh /app/entrypoint.sh
