FROM golang:1.10 as builder
WORKDIR /go/src/github.com/scofieldpeng/dnspod-ddns/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scofieldpeng/alpine:glibc-2.7
RUN mkdir /app
COPY --from=builder /go/src/github.com/scofieldpeng/dnspod-ddns/app /app/app
RUN chmod +x /app/app
WORKDIR /app/

ENTRYPOINT /app/app