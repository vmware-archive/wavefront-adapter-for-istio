FROM golang:1.11 as builder
MAINTAINER Venil Noronha <veniln@vmware.com>

WORKDIR /go/src/github.com/vmware/wavefront-adapter-for-istio/
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/wavefront ./wavefront/cmd/

FROM photon:2.0
WORKDIR /bin/
COPY --from=builder /go/src/github.com/vmware/wavefront-adapter-for-istio/bin/wavefront .
ENTRYPOINT [ "/bin/wavefront" ]
CMD [ "8080" ]
EXPOSE 8080
