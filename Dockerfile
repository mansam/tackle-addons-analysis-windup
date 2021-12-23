FROM registry.access.redhat.com/ubi8/go-toolset:1.16.7 as addon-adapter-builder
ENV GOPATH=$APP_ROOT
COPY go.mod go.sum main.go Makefile .
RUN make addon-adapter

FROM registry.access.redhat.com/ubi8/openjdk-8-runtime:latest
COPY mta-cli /opt/mta-cli
COPY --from=addon-adapter-builder /opt/app-root/src/bin/addon-adapter /usr/local/bin/addon-adapter

ENTRYPOINT ["/usr/local/bin/addon-adapter"]

