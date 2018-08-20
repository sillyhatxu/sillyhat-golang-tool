FROM golang:1.10
LABEL MAINTAINER="heixiushamao@gmail.com"


RUN go get -u github.com/bsm/sarama-cluster


ENV ENVIROMENT_PROFILES_ACTIVE dt
ENV TIME_ZONE=Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone
RUN apk add --no-cache tzdata

##CA证书，https请求
RUN apk add --no-cache ca-certificates

WORKDIR /go
COPY main /go
COPY . /go/payment

##CMD ["./main","-enviroment=dt"]
##CMD ["./main","-enviroment=dp"]
##CMD ["./main","-enviroment=production"]
ENTRYPOINT ./main -enviroment=$ENVIROMENT_PROFILES_ACTIVE