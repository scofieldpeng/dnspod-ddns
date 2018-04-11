package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type appConfig struct {
	dnspodId    string
	dnspodToken string
	recordId    string
	domain      string
	subDomain   string
	internal    int
	email       string
}

// 验证appConfig
func (c appConfig) Validate() (err error) {
	if c.dnspodId == "" {
		return errors.New("environment DNSPOD_ID required")
	}
	if c.dnspodToken == "" {
		return errors.New("environment DNSPOD_TOKEN required")
	}

	if c.recordId == "" && c.subDomain == "" {
		return errors.New("environment DNSPOD_RECORDID or DNSPOD_SUBDOMAIN required")
	}
	if c.domain == "" {
		return errors.New("environment DNSPOD_DOMAIN required")
	}
	if c.internal < 5 {
		return errors.New("environment DNSPOD_INTERNAL should range from 0 to 5")
	}

	return
}

var (
	config = appConfig{
		dnspodId:    os.Getenv("DNSPOD_ID"),
		dnspodToken: os.Getenv("DNSPOD_TOKEN"),
		domain:      os.Getenv("DNSPOD_DOMAIN"),
		subDomain:   os.Getenv("DNSPOD_SUBDOMAIN"),
		internal:    60,
		email:       os.Getenv("DNSPOD_EMAIL"),
	}
)

const (
	ClientUserAgent = "DNSPOD-DDNS-CLIENT"
	Version         = "1.0.0"
	StatusOk        = "1"
)

func init() {
	internal := os.Getenv("DNSPOD_INTERNAL")
	config.internal, _ = strconv.Atoi(internal)
	if config.internal < 5 {
		config.internal = 60
	}
	if config.email == "" {
		config.email = "example@example.com"
	}
}

func main() {
	var (
		err          error
		lastPublicIP string
		publicIP     string
	)

	if err := config.Validate(); err != nil {
		fmt.Println("[error]", err)
		os.Exit(1)
	}

	fmt.Println("start")
	for {
		publicIP, err = GetPublicIP()
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Duration(config.internal) * time.Second)
			continue
		}
		if config.recordId == "" || lastPublicIP == "" {
			config.recordId, lastPublicIP, err = GetRecord()
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Duration(config.internal) * time.Second)
				continue
			}
		}
		if publicIP != lastPublicIP {
			fmt.Println("发现公网IP变化，开始更新")
			if err = UpdateRecord(config.recordId, publicIP); err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Duration(config.internal) * time.Second)
				continue
			}
			fmt.Println("公网IP更新成功，新的公网IP:", publicIP)
			lastPublicIP = publicIP
		}
		fmt.Println("下次更新时间:", time.Now().Add(time.Duration(config.internal)*time.Second).Format("2006-01-02 15:04:05"))
		time.Sleep(time.Duration(config.internal) * time.Second)
	}
}

// 公共返回参数
type CommonResponse struct {
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		CreateTime string `json:"created_at"`
	} `json:"status"`
}

// 记录列表返回值
type RecordListResponse struct {
	CommonResponse
	Records []struct {
		SubDomain string `json:"name"`
		Id        string `json:"id"`
		PublicIP  string `json:"value"`
	} `json:"records"`
}

// 更新record记录
func UpdateRecord(recordId string, publicIP string) (err error) {
	var (
		request      *http.Request
		response     *http.Response
		c            *http.Client
		body         = url.Values{}
		responseData CommonResponse
	)
	body.Add("login_token", fmt.Sprintf("%s,%s", config.dnspodId, config.dnspodToken))
	body.Add("format", "json")
	body.Add("lang", "cn")
	body.Add("error_on_empty", "no")
	body.Add("domain", config.domain)
	body.Add("sub_domain", config.subDomain)
	body.Add("record_id", recordId)
	body.Add("record_type", "A")
	body.Add("record_line", "默认")
	body.Add("value", publicIP)

	request, err = http.NewRequest("POST", "https://dnsapi.cn/Record.Modify", strings.NewReader(body.Encode()))
	if err != nil {
		err = errors.New("request对象创建失败,err:" + err.Error())
		return
	}
	request.Header.Add("User-Agent", fmt.Sprintf("%s/%s(%s)", ClientUserAgent, Version, config.email))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c = &http.Client{Timeout: time.Second * 30}
	response, err = c.Do(request)
	if err != nil {
		err = errors.New("请求出错,err:" + err.Error())
		return
	}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		err = errors.New("解析数据失败，err:" + err.Error())
		return
	}
	defer response.Body.Close()

	if responseData.Status.Code != StatusOk {
		err = errors.New(fmt.Sprintf("更新失败,code:%s,message:%s", responseData.Status.Code, responseData.Status.Message))
		return
	}

	return
}

// 获取recordid
func GetRecord() (recordId, IP string, err error) {
	var (
		request      *http.Request
		response     *http.Response
		c            *http.Client
		body         = url.Values{}
		responseData RecordListResponse
	)
	body.Add("login_token", fmt.Sprintf("%s,%s", config.dnspodId, config.dnspodToken))
	body.Add("format", "json")
	body.Add("lang", "cn")
	body.Add("error_on_empty", "no")

	body.Add("domain", config.domain)
	body.Add("sub_domain", config.subDomain)
	request, err = http.NewRequest("POST", "https://dnsapi.cn/Record.List", strings.NewReader(body.Encode()))
	if err != nil {
		err = errors.New("request对象创建失败,err:" + err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", fmt.Sprintf("%s/%s(%s)", ClientUserAgent, Version, config.email))
	c = &http.Client{Timeout: time.Second * 30}
	response, err = c.Do(request)
	if err != nil {
		err = errors.New("请求出错,err:" + err.Error())
		return
	}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		err = errors.New("解析数据失败，err:" + err.Error())
		return
	}
	defer response.Body.Close()
	if responseData.Status.Code != StatusOk {
		err = errors.New(fmt.Sprintf("获取record失败,code:%s,message:%s", responseData.Status.Code, responseData.Status.Message))
		return
	}

	for _, v := range responseData.Records {
		if v.SubDomain == config.subDomain {
			recordId = v.Id
			IP = v.PublicIP
			return
		}
	}

	err = errors.New("没有找到相关记录，请先前往dnspod进行添加")
	return
}

type getPublicIPResponse struct {
	IP string `json:"origin"`
}

// 获取公网IP,如果出错，返回第二个参数
func GetPublicIP() (publicIP string, err error) {
	var (
		response     *http.Response
		responseData getPublicIPResponse
	)
	if response, err = http.Get("http://www.httpbin.org/ip"); err != nil {
		err = errors.New("获取公网IP出错,err:" + err.Error())
		return
	}
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		err = errors.New("获取公网IP出错,err:" + err.Error())
		return
	}
	defer response.Body.Close()

	return responseData.IP, nil
}
