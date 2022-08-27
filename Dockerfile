FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

RUN mkdir /etc/asoul-video
WORKDIR /etc/asoul-video

ADD asoul-video /etc/asoul-video

RUN chmod 655 /etc/asoul-video/asoul-video

ENTRYPOINT ["/etc/asoul-video/asoul-video"]
EXPOSE 2830
