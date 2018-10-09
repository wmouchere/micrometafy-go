FROM golang:latest

WORKDIR /go/src/github.com/wmouchere/micrometafy-go/micrometafyplaylist-go
COPY . .

RUN go install -v

EXPOSE 8080

CMD ["micrometafyplaylist-go"]