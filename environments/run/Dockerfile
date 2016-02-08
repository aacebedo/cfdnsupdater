from gliderlabs/alpine

RUN apk add  --no-cache ca-certificates

ARG CFDNSUPDATER_VERSION
ENV CFDNSUPDATER_VERSION ${CFDNSUPDATER_VERSION:-0.6}

ARG CFDNSUPDATER_ARCH
ENV CFDNSUPDATER_ARCH ${CFDNSUPDATER_ARCH:-amd64}


ADD https://github.com/aacebedo/cfdnsupdater/releases/download/${CFDNSUPDATER_VERSION}/cfdnsupdater.${CFDNSUPDATER_ARCH}.${CFDNSUPDATER_VERSION}.tar.gz /tmp/cfdnsupdater.tar.gz
RUN tar xvzf /tmp/cfdnsupdater.tar.gz -C /usr/local/bin
RUN rm -rf /tm p/cfdnsupdater.tar.gz

ENTRYPOINT ["cfdnsupdater"]

