FROM golang:latest
COPY . /untitled_api
WORKDIR /untitled_api

#RUN go get magin.go
#RUN GOOS=linux GOARCH=amd64 go build main.go

CMD ["/untitled_api/main"]

EXPOSE 8081
