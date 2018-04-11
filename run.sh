#!/usr/bin/env bash

# 填写你在dnspod申请的id
DNSPOD_ID=123456
# 在dnspod申请的token
DNSPOD_TOKEN=123456
# 在dnspod要更新的顶级域名
DNSPOD_DOMAIN=example.com
# 在dnspod要更新的子域名前缀，如果是根域名，填写@即可
DNSPOD_SUBDOMAIN=example
# 你的邮箱
DNSPOD_EMAIL=example@example.com

docker run --name=ddns --restart=always -d \
    -e DNSPOD_ID=${DNSPOD_ID} \
    -e DNSPOD_TOKEN=${DNSPOD_TOKEN} \
    -e DNSPOD_DOMAIN=${DNSPOD_DOMAIN} \
    -e DNSPOD_SUBDOMAIN=${DNSPOD_SUBDOMAIN} \
    -e DNSPOD_EMAIL=${DNSPOD_EMAIL} \
    scofieldpeng/dnspod-ddns:1.0.0