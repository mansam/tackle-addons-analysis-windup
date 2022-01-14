FROM registry.access.redhat.com/ubi8/go-toolset:1.16.7 as builder
ENV GOPATH=$APP_ROOT
COPY . .
RUN make addon

FROM registry.access.redhat.com/ubi8/openjdk-8-runtime:latest
COPY mta-cli /opt/mta-cli
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]
