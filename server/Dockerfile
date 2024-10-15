FROM golang:1.23@sha256:a7f2fc9834049c1f5df787690026a53738e55fc097cd8a4a93faa3e06c67ee32 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/beansd

FROM alpine@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

RUN mkdir /app
COPY --from=buildgo /app/beansd /app/

ENV BEANS_HTTP_PORT=8080

CMD ["/app/beansd"]
