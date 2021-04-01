# FROM golang:1.15 as builder
FROM docker.dm-ai.cn/public/golang:1.15-nvidia-gpu-mem-monitor as builder
WORKDIR /data
ADD . .
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go build -o targets/nvidia-gpu-mem-monitor

# FROM alpine:3.12.0
FROM nvidia/cuda:11.0-base
WORKDIR /data
COPY --from=builder /data/targets/nvidia-gpu-mem-monitor .
CMD /data/nvidia-gpu-mem-monitor
EXPOSE 80