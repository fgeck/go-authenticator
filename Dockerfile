FROM golang:1.17.8-alpine as builder

WORKDIR /

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /authenticator main.go

FROM scratch

COPY --from=builder /authenticator /authenticator
#RUN ls -lah
EXPOSE 9123

ENTRYPOINT ["/authenticator"]
