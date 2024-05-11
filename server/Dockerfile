FROM golang:1.22@sha256:b1e05e2c918f52c59d39ce7d5844f73b2f4511f7734add8bb98c9ecdd4443365 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/beansd

FROM alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

RUN mkdir /app
COPY --from=buildgo /app/beansd /app/

ENV BEANS_HTTP_PORT=8080

CMD ["/app/beansd"]
