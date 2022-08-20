FROM golang:1.18 as builder
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o server ./main.go

FROM alpine
COPY --from=builder /app/server /app/server
EXPOSE 8000
CMD ["/app/server"]
