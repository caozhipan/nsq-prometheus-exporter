# nsq-prometheus-exporter

## BUILD
```bash
go build -o nsq-prometheus-exporter main.go
```

## DOCKER
```bash
docker run -p 9527:9527  nsq-prometheus-exporter -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```

## RUN

```bash
./nsq-prometheus-exporter -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```