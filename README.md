# dnspod-ddns

通过dnspod来进行ddns的小脚本，用来做什么，你懂得，本来不想写的，只是github.com另外一个哥们儿用py写的跑不起来，看了下代码有问题，懒得修改直接用golang撸了一个

如果对你有用，可以star一下，2333

## 使用

**配置项:**

配置采用环境变量驱动，具体见下：

```ini
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
```

### docker运行

```bash
docker run --name=ddns --restart=always -d \
    -e DNSPOD_ID=${DNSPOD_ID} \
    -e DNSPOD_TOKEN=${DNSPOD_TOKEN} \
    -e DNSPOD_DOMAIN=${DNSPOD_DOMAIN} \
    -e DNSPOD_SUBDOMAIN=${DNSPOD_SUBDOMAIN} \
    -e DNSPOD_EMAIL=example@example.com \
    scofieldpeng/dnspod-ddns:1.0.0
```

### 源码编译

1. 安装go
2. 下载源码
```
go get github.com/scofieldpeng/dnspod-ddns
```
3. 编译
```
go build -o app .
```

$GOPATH/src/github.com/scofieldpeng/dnspod-ddns/app文件即为二进制包，直接运行即可

## 申请dnspod的ID和token

教程详见dnspod官网[https://support.dnspod.cn/Kb/showarticle/tsid/227/](https://support.dnspod.cn/Kb/showarticle/tsid/227/)
