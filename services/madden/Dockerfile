FROM docker-hfrd.di2e.net/golang:1.16.3-alpine3.13 AS build

ARG USER=user
ARG PASS=pass

RUN apk add build-base git
ENV GOPRIVATE=bitbucket.di2e.net
RUN printf "machine bitbucket.di2e.net\n\
    login ${USER}\n\
    password ${PASS}\n\
    \n"\
    >> /root/.netrc
RUN chmod 600 /root/.netrc
RUN cat /root/.netrc

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download

COPY swagger ./swagger
COPY dataservice ./dataservice
COPY controller ./controller
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o /madden .

FROM scratch AS runtime
COPY --from=build /madden ./
EXPOSE 8080/tcp
ENTRYPOINT ["./madden"]
