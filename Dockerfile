#使用Golang1.22.5镜像作为基础镜像
FROM golang:1.22.5
#设置工作目录
WORKDIR /app
#将当前目录的所有文件复制到工作目录
COPY . .
#下载并安装依赖
RUN go mod tidy
#编译Go应用
RUN go build -o main .
#暴露应用运行端口1323
EXPOSE 1323
#运行应用
CMD ["./main"]
