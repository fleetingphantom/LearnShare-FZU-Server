FROM golang:1.25.4 AS builder

ENV TZ=Asia/Shanghai
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

RUN mkdir -p /app
WORKDIR /app
ADD . /app
RUN go mod tidy
RUN sh build.sh

FROM alpine
RUN apk update --no-cache && \
    apk add --no-cache ca-certificates tzdata ffmpeg

ENV TZ=Asia/Shanghai
ENV service=LearnShare

WORKDIR /app
COPY --from=builder /app/output /app
COPY --from=builder /app/script/bootstrap.sh /app/
COPY --from=builder /app/config/config.example.yaml /app/config/

CMD ["sh","bootstrap.sh"]
