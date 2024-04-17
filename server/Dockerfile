FROM golang:1.22@sha256:450e3822c7a135e1463cd83e51c8e2eb03b86a02113c89424e6f0f8344bb4168 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/beansd

FROM alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

RUN mkdir /app
COPY --from=buildgo /app/beansd /app/

ENV BEANS_HTTP_PORT=8080

CMD ["/app/beansd"]
