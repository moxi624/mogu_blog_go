FROM golang as builder
WORKDIR /go/src/mogu_blog_go
COPY . .
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.io,direct && go install github.com/beego/bee && bee pack -be GOOS=linux

FROM ubuntu:latest
ENV LANG C.UTF-8
ENV TZ Asia/Shanghai
WORKDIR /go/src/mogu_blog_go
COPY --from=builder /go/src/mogu_blog_go/mogu_blog_go.tar.gz .
RUN  tar -zxvf mogu_blog_go.tar.gz
RUN  rm -rf mogu_blog_go.tar.gz
RUN  chmod +x mogu_blog_go
EXPOSE 8607
ENTRYPOINT ./mogu_blog_go