from alpine
RUN apk add --update go git python3 bash

ENV GOPATH /go:/out/cfdnsupdater
RUN go get gopkg.in/alecthomas/kingpin.v2
RUN go get github.com/op/go-logging
RUN go get github.com/levigross/grequests
ENV CGO_ENABLED=0

ADD ./build.py /usr/bin
RUN chmod u+x /usr/bin/build.py

WORKDIR /out
VOLUME ["/out"]

ENTRYPOINT ["build.py"]



