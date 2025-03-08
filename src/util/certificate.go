package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"pm-ssl-management/src/conf"
	"pm-ssl-management/src/log"
	"strings"
	"time"

	cas20200407 "github.com/alibabacloud-go/cas-20200407/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	goutils "github.com/typa01/go-utils"
)

func CreateClient(Endpoint string) (_result *cas20200407.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(conf.GetConf("AccessKeyId")),
		AccessKeySecret: tea.String(conf.GetConf("AccessKeySecret")),
	}
	config.Endpoint = tea.String(Endpoint)
	_result = &cas20200407.Client{}
	_result, _err = cas20200407.NewClient(config)
	return _result, _err
}

func CreateApiInfo(name string, Version string) (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String(name),
		// 接口版本
		Version: tea.String(Version),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

func DeployCertificate(cert string, key string, certPath string, keyPath string, command string, domain string) {
	// 创建证书
	guid := goutils.GUID()
	certificateTemPath := path.Join(GetWdFile("tmp"), guid)
	// 创建临时证书文件夹
	os.MkdirAll(certificateTemPath, 0755)
	// 创建证书
	certTmpPath := path.Join(certificateTemPath, "cert.pem")
	certFile, _ := os.Create(certTmpPath)
	certFile.WriteString(cert)
	certFile.Close()
	// 创建私钥
	keyTmpPath := path.Join(certificateTemPath, "key.pem")
	keyFile, _ := os.Create(keyTmpPath)
	keyFile.WriteString(key)
	keyFile.Close()
	// 备份原始证书
	formattedTime := time.Now().Format("20060102_150405") // 备份的时间戳
	certPath = filepath.ToSlash(certPath)                 // 将路径分隔符统一
	keyPath = filepath.ToSlash(keyPath)                   // 将路径分隔符统一
	// 是否存在原始证书，有就先备份
	if _, err := os.Stat(certPath); err == nil {
		os.Rename(certPath, fmt.Sprintf("%s.bak_%s", certPath, formattedTime))
	}
	if _, err := os.Stat(keyPath); err == nil {
		os.Rename(keyPath, fmt.Sprintf("%s.bak_%s", keyPath, formattedTime))
	}
	// 拷贝完将临时证书删除
	defer func() {
		os.RemoveAll(certificateTemPath)
	}()

	// 先判断是否有存放证书的文件夹，没有就需要创建
	folderCertPath := filepath.Dir(certPath)
	if _, err := os.Stat(folderCertPath); os.IsNotExist(err) {
		os.MkdirAll(folderCertPath, 0755)
	}
	if err := os.Rename(certTmpPath, certPath); err != nil {
		log.ComLoggerClient.Error(fmt.Sprintf("域名 %s 私钥移动失败: %s", domain, err.Error()))
	}
	if err := os.Rename(keyTmpPath, keyPath); err != nil {
		log.ComLoggerClient.Error(fmt.Sprintf("域名 %s 私钥移动失败: %s", domain, err.Error()))
	}

	// 执行重启命令
	if outPut, err := RestartService(command); err != nil {
		log.ComLoggerClient.Error(fmt.Sprintf("域名 %s 证书更新命令执行失败: %s", domain, err.Error()))
	} else {
		log.ComLoggerClient.Info(fmt.Sprintf("域名 %s 证书更新命令执行成功: %s", domain, outPut))
	}
}

func RestartService(commandStr string) (out string, err error) {
	cmdArgs := strings.Fields(commandStr)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func CanApplyCertificate() error {

	client, _err := CreateClient("cas.aliyuncs.com")
	if _err != nil {
		return _err
	}

	params := CreateApiInfo("DescribePackageState", "2020-04-07")
	queries := map[string]interface{}{}
	queries["ProductCode"] = tea.String("digicert-free-1-free")
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	data, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return _err
	}
	body, _ := data["body"].(map[string]interface{})

	log.ComLoggerClient.Info("DescribePackageState: ", body)

	usedCount, err := body["UsedCount"].(json.Number).Int64()
	if err != nil {
		return err
	}

	totalCount, err := body["TotalCount"].(json.Number).Int64()
	if err != nil {
		return err
	}

	if usedCount >= totalCount {
		return errors.New("证书申请次数已用完")
	}

	return nil
}

func CreateCertificateForPackageRequest(domain string) (int64, error) {
	client, _err := CreateClient("cas.aliyuncs.com")
	if _err != nil {
		return 0, _err
	}
	params := CreateApiInfo("CreateCertificateForPackageRequest", "2020-04-07")
	queries := map[string]interface{}{}
	queries["ProductCode"] = tea.String("digicert-free-1-free")
	queries["Domain"] = tea.String(domain)
	queries["ValidateType"] = tea.String("DNS")
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	data, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return 0, _err
	}
	body := data["body"].(map[string]interface{})

	log.ComLoggerClient.Info("CreateCertificateForPackageRequest: ", body)

	OrderId, _ := body["OrderId"].(json.Number).Int64()

	return OrderId, nil
}

func DescribeCertificateState(OrderId int) (string, string, string, error) {
	client, _err := CreateClient("cas.aliyuncs.com")
	if _err != nil {
		return "", "", "", _err
	}

	params := CreateApiInfo("DescribeCertificateState", "2020-04-07")
	queries := map[string]interface{}{}
	queries["OrderId"] = tea.Int(OrderId)
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	data, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return "", "", "", _err
	}
	body := data["body"].(map[string]interface{})

	log.ComLoggerClient.Info("DescribeCertificateState: ", body)

	if body["Type"].(string) == "certificate" {
		return body["Type"].(string), body["Certificate"].(string), body["PrivateKey"].(string), nil
	}

	if body["Type"].(string) == "domain_verify" {
		return body["Type"].(string), body["RecordDomain"].(string), body["RecordValue"].(string), nil
	}

	return body["Type"].(string), "", "", nil
}

func GetCertificateTime(orderId int64) (time.Time, time.Time, error) {
	var startTime = time.Now()
	var endTime = time.Now().Add(time.Hour * 24 * 90)
	// 先查询全部证书
	client, _err := CreateClient("cas.aliyuncs.com")
	if _err != nil {
		return startTime, endTime, _err
	}

	params := CreateApiInfo("ListUserCertificateOrder", "2020-04-07")
	// query params
	queries := map[string]interface{}{}
	queries["ShowSize"] = tea.Int(1000)
	// runtime options
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值实际为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	data, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return startTime, endTime, _err
	}

	body := data["body"].(map[string]interface{})

	log.ComLoggerClient.Info("DescribeCertificateState: ", body)

	certificateList := body["CertificateOrderList"].([]interface{})

	for _, v := range certificateList {
		v := v.(map[string]interface{})
		OrderIdV, _ := v["OrderId"].(json.Number).Int64()
		CertStartTime, _ := v["CertStartTime"].(json.Number).Int64()
		CertEndTime, _ := v["CertEndTime"].(json.Number).Int64()
		if OrderIdV == orderId {
			startTime := time.Unix(CertStartTime/1000, 0)
			endTime := time.Unix(CertEndTime/1000, 0)
			return startTime, endTime, nil
		}
	}
	return startTime, endTime, errors.New("未找到对应证书")
}

func CreateDnsRecord(domain string, RR string, value string) error {
	client, _err := CreateClient("alidns.cn-hangzhou.aliyuncs.com")
	if _err != nil {
		return _err
	}

	// 按照 "." 分割域名
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	params := CreateApiInfo("AddDomainRecord", "2015-01-09")
	queries := map[string]interface{}{}
	queries["DomainName"] = tea.String(domain)
	queries["RR"] = tea.String(strings.Replace(RR, "."+domain, "", 1))
	queries["Type"] = tea.String("TXT")
	queries["Value"] = tea.String(value)
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	_, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return _err
	}
	return nil
}

func GetUserCertificateOrder() {
	client, err := CreateClient("cas.aliyuncs.com")
	if err != nil {
		log.ComLoggerClient.Error("GetUserCertificateOrder: ", err.Error())
		return
	}
	params := CreateApiInfo("ListUserCertificateOrder", "2020-04-07")
	queries := map[string]interface{}{}
	queries["ShowSize"] = tea.Int(1000)
	runtime := &service.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	data, err := client.CallApi(params, request, runtime)
	if err != nil {
		log.ComLoggerClient.Error("GetUserCertificateOrder: ", err.Error())
		return
	}
	body := data["body"].(map[string]interface{})
	log.ComLoggerClient.Info("GetUserCertificateOrder: ", body)
}
