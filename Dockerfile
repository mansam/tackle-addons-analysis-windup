FROM registry.access.redhat.com/ubi8/go-toolset:1.16.7 as builder
ENV GOPATH=$APP_ROOT
COPY . .
RUN make addon

FROM registry.access.redhat.com/ubi8/openjdk-8-runtime
USER 0
RUN echo -e "[centos8]" \
      "\nname = centos8" \
      "\nbaseurl = http://mirror.centos.org/centos/8-stream/AppStream/x86_64/os/" \
      "\nenabled = 1" \
      "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN microdnf -y install git subversion \
 && microdnf -y clean all
USER 185
CMD wget -O  /opt https://repo1.maven.org/maven2/org/jboss/windup/mta-cli/5.2.1.Final/mta-cli-5.2.1.Final-offline.zip \
    && unzip /opt/mta-cli-5.2.1.Final-offline.zip -d /opt \
    && rm /opt/mta-cli-5.2.1.Final-offline.zip
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]
