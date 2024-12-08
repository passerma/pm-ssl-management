package controller

import (
	"pm-ssl-management/src/conf"
	"pm-ssl-management/src/cron"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/model"
	"pm-ssl-management/src/util"
	"strconv"

	cas20200407 "github.com/alibabacloud-go/cas-20200407/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/gin-gonic/gin"
)

type CertificateList struct {
	CommonName    string
	EndDate       string
	StartDate     string
	CertificateId int64
}

func CreateClient() (_result *cas20200407.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(conf.GetConf("AccessKeyId")),
		AccessKeySecret: tea.String(conf.GetConf("AccessKeySecret")),
	}
	config.Endpoint = tea.String("cas.aliyuncs.com")
	_result = &cas20200407.Client{}
	_result, _err = cas20200407.NewClient(config)
	return _result, _err
}

func CreateApiInfo(name string) (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String(name),
		// 接口版本
		Version: tea.String("2020-04-07"),
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

func PostCertificate(ctx *gin.Context) {
	var req map[string]interface{}
	ctx.BindJSON(&req)

	// 校验参数
	if req["domain"] == nil || req["keyPath"] == nil || req["certPath"] == nil ||
		req["command"] == nil || req["autoRenew"] == nil || req["renewTime"] == nil {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("保存证书信息失败: 参数错误")
		return
	}

	var RenewTime float64 = 7
	if v := req["renewTime"].(float64); v > 0 && v < 366 {
		RenewTime = v
	}

	certificate := model.Certificate{
		Domain:    req["domain"].(string),
		KeyPath:   req["keyPath"].(string),
		CertPath:  req["certPath"].(string),
		Command:   req["command"].(string),
		AutoRenew: req["autoRenew"].(bool),
		RenewTime: int64(RenewTime),
	}

	if err := model.DdClient.Create(&certificate).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(3))
		log.ComLoggerClient.Error("保存证书信息失败: ", err.Error())
		return
	}

	ctx.JSON(200, util.SendSusModel(certificate))
}

func GetCertificate(ctx *gin.Context) {
	var certificateList []model.Certificate
	model.DdClient.Find(&certificateList)
	ctx.JSON(200, util.SendSusModel(certificateList))
}

func DeleteCertificate(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("删除证书信息失败: 参数错误")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("删除证书信息失败: 参数错误")
		return
	}
	if err := model.DdClient.Delete(&model.Certificate{}, idInt).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(3))
		log.ComLoggerClient.Error("删除证书信息失败: ", err.Error())
	}
	ctx.JSON(200, util.SendSusModel(nil))
}

func PutCertificate(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("更新证书信息失败: 参数错误")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("更新证书信息失败: 参数错误")
		return
	}
	var certificate model.Certificate
	model.DdClient.First(&certificate, idInt)

	if certificate.ID == 0 {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("更新证书信息失败: 未申请对应证书")
		return
	}

	var req map[string]interface{}
	ctx.BindJSON(&req)

	if req["keyPath"] != nil {
		certificate.KeyPath = req["keyPath"].(string)
	}
	if req["certPath"] != nil {
		certificate.CertPath = req["certPath"].(string)
	}
	if req["command"] != nil {
		certificate.Command = req["command"].(string)
	}
	if req["autoRenew"] != nil {
		certificate.AutoRenew = req["autoRenew"].(bool)
	}
	if req["renewTime"] != nil {
		var RenewTime float64 = 7
		if v := req["renewTime"].(float64); v > 0 && v < 366 {
			RenewTime = v
		}
		certificate.RenewTime = int64(RenewTime)
	}

	if err := model.DdClient.Save(&certificate).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(3))
		log.ComLoggerClient.Error("更新证书信息失败: ", err.Error())
		return
	}
	ctx.JSON(200, util.SendSusModel(nil))
}

func GetCertificateState(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("获取证书状态失败: 参数错误")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("获取证书状态失败: 参数错误")
		return
	}

	var certificate model.Certificate
	if err := model.DdClient.First(&certificate, idInt).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(1, "未申请对应证书"))
		log.ComLoggerClient.Error("获取证书状态失败: ", err.Error())
		return

	}
	if certificate.OrderId == nil || *certificate.OrderId == 0 {
		ctx.JSON(200, util.SendErrModel(1, "未申请对应证书"))
		log.ComLoggerClient.Error("获取证书状态失败: 未申请对应证书")
		return
	}
	if state, _, _, err := util.DescribeCertificateState(int(*certificate.OrderId)); err != nil {
		ctx.JSON(200, util.SendErrModel(1, err.Error()))
		log.ComLoggerClient.Error("获取证书状态失败: ", err.Error())
		return
	} else {
		ctx.JSON(200, util.SendSusModel(state))
	}
}

// {功能} 申请证书
// {参数} 路由
// {返回} 无
func PostCertificateApply(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("申请证书失败: 参数错误")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, util.SendErrModel(2))
		log.ComLoggerClient.Error("申请证书失败: 参数错误")
		return
	}
	var certificate model.Certificate
	if err := model.DdClient.First(&certificate, idInt).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(2, "未找到对应证书"))
		log.ComLoggerClient.Error("申请证书失败: ", "未找到对应证书")
		return
	}

	if err := util.CanApplyCertificate(); err != nil {
		ctx.JSON(200, util.SendErrModel(1, err.Error()))
		log.ComLoggerClient.Error("申请证书失败: ", err.Error())
		return
	}

	if orderId, err := util.CreateCertificateForPackageRequest(certificate.Domain); err != nil {
		ctx.JSON(200, util.SendErrModel(1, err.Error()))
		log.ComLoggerClient.Error("申请证书失败: ", err.Error())
		return
	} else {
		certificate.OrderId = &orderId
	}

	certificate.State = "domain_verify"
	certificate.Validate = false

	if err := model.DdClient.Save(&certificate).Error; err != nil {
		ctx.JSON(200, util.SendErrModel(1))
		log.ComLoggerClient.Error("申请证书失败: ", err.Error())
		return
	}

	// 执行定时任务，检查证书状态
	cron.AddCertificateState(certificate.ID, *certificate.OrderId)

	ctx.JSON(200, util.SendSusModel(nil))
}
