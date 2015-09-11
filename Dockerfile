FROM golang

ADD . /go/src/github.com/tcyrus/hackaday-io-badges
RUN go install github.com/tcyrus/hackaday-io-badges
ENTRYPOINT /go/bin/hackaday-io-badges

EXPOSE 8080
