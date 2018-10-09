FROM golang:latest

WORKDIR /go/src/github.com/wmouchere/micrometafy-go/micrometafyquery-go
COPY . .

RUN go install -v

EXPOSE 8080

CMD ["micrometafyquery-go"]