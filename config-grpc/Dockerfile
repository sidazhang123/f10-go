FROM alpine
ENV CONFIG_GRPC_ADDR :9600
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.10/main/" > /etc/apk/repositories
RUN apk update && apk add tzdata && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD config-grpc /config-grpc
RUN mkdir conf
ENTRYPOINT [ "/config-grpc" ]