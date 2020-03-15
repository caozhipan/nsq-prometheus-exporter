# nsq-prometheus-exporter

## BUILD
```bash
go build -o nsq-prometheus-exporter main.go
```

## DOCKER

- docker build & docker run
```bash
docker build -t nsq-prometheus-exporter .
docker run -p 9527:9527  nsq-prometheus-exporter -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```

- use an official image
```bash
docker run -p 9527:9527  whoisyourdady/nsq-prometheus-exporter:latest -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```

## RUN

```bash
./nsq-prometheus-exporter -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```
