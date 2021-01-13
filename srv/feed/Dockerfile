FROM alpine

ENV CONFIG_GRPC_ADDR 127.0.0.1:9600
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.10/main/" > /etc/apk/repositories
RUN apk update && apk add tzdata && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD feed-srv /feed-srv

CMD [ "/feed-srv","--broker=http" ,"--broker_address=0.0.0.0:0" ]