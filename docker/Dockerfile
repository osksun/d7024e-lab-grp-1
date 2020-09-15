FROM ubuntu:18.04

WORKDIR /app

RUN apt-get update && apt-get install -y fping
RUN apt-get update && apt-get install -y golang-go

ADD pingAll.sh /app
ADD net1 /app

#COPY . .

RUN go build -o main .
#RUN go build
#RUN export GOPATH=$HOME
#RUN go run net.go
#RUN go install example.com/user/net
#RUN net

#CMD ["go","run","net.go"]
#EXPOSE 8080

ENTRYPOINT ["/app/main"]
