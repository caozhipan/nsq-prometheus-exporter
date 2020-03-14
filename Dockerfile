FROM golang:1.14-alpine as builder

WORKDIR /workspace
ADD . ./

RUN CGO_ENABLED=0  go build -ldflags "-w -s" -i -o nsq-prometheus-exporter main.go

FROM scratch
WORKDIR /workspace
COPY --from=builder /workspace/nsq-prometheus-exporter ./
EXPOSE 9527
ENTRYPOINT ["./nsq-prometheus-exporter"]
