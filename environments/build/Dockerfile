from alpine
RUN apk add --update go git python3 bash

ENV GOPATH /go
ENV CGO_ENABLED=0

ADD ./build.py /usr/bin
RUN chmod u+x /usr/bin/build.py

WORKDIR /out
VOLUME ["/out"]

ENTRYPOINT ["build.py", "-p", "aacebedo/cfdnsupdater", "-n", "cfdnsupdater"]



