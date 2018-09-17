FROM photon:2.0 as builder
MAINTAINER Venil Noronha <veniln@vmware.com>

RUN tdnf update -y
RUN tdnf install -y go-1.10.3-1.ph2.x86_64
WORKDIR /root/go/src/github.com/vmware/wavefront-adapter-for-istio/
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/wavefront ./wavefront/cmd/

FROM photon:2.0
RUN tdnf update -y
WORKDIR /bin/
COPY --from=builder /root/go/src/github.com/vmware/wavefront-adapter-for-istio/bin/wavefront .
COPY open_source_license .
ENTRYPOINT [ "/bin/wavefront" ]
CMD [ "8000" ]
EXPOSE 8000
