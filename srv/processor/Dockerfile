FROM golang:1.13.5-alpine
RUN apk update && apk add gcc libc-dev make git
ENV CONFIG_GRPC_ADDR 127.0.0.1:9600
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.10/main/" > /etc/apk/repositories
RUN apk add tzdata && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /plugins/src
RUN mkdir -p /plugins/so
COPY plugins/src/ /plugins/src/
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go get github.com/olekukonko/tablewriter
RUN cd /plugins/src && go build -buildmode=plugin -o /plugins/so/latest_tips.so latest_tips.go && \
go build -buildmode=plugin -o /plugins/so/financial_analysis.so financial_analysis.go && \
go build -buildmode=plugin -o /plugins/so/operational_analysis.so operational_analysis.go && \
go build -buildmode=plugin -o /plugins/so/shareholder_analysis.so shareholder_analysis.go
ADD processor-srv /processor-srv
WORKDIR /
ENTRYPOINT [ "/processor-srv" ]
