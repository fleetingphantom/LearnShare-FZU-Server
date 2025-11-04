FROM golang:1.23.6 AS builder

ENV TZ Asia/Shanghai
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

RUN mkdir -p /app

WORKDIR /app

ADD docker /app
RUN go mod tidy
RUN sh build.sh

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata ffmpeg
ENV TZ Asia/Shanghai
ENV service LearnShare

WORKDIR /app

COPY --from=builder /app/output /app/output
ADD ./docker/bootstrap.sh /app/
EXPOSE 8888
CMD ["sh","bootstrap.sh"]
