package cron

import (
	"fmt"
	"pm-ssl-management/src/log"
	"pm-ssl-management/src/model"
	"pm-ssl-management/src/util"
	"time"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

var certificateStateMap = make(map[uint]int64)

func initCertificateStateMap() {
	var certificates []model.Certificate
	// 只有状态不是certificate或者verify_fail的证书才需要监听
	if err := model.DdClient.
		Where("state != ?", "certificate").
		Where("state != ?", "verify_fail").
		Where("order_id IS NOT NULL").
		Find(&certificates).Error; err != nil {
		log.ComLoggerClient.Error("获取监听证书状态失败: ", err.Error())
	} else {
		log.ComLoggerClient.Info("获取监听证书状态成功，证书数据: ", certificates)
	}
	for _, v := range certificates {
		certificateStateMap[v.ID] = *v.OrderId
	}
}

func initCertificateStateCron() {
	go initCertificateStateMap()
	c := cron.New()
	c.AddFunc("*/5 * * * * *", func() {
		for k, v := range certificateStateMap {
			log.ComLoggerClient.Info(fmt.Sprintf("获取证书状态: id: %d, 订单: %d", k, v))
			if state, cert, key, err := util.DescribeCertificateState(int(v)); err != nil {
				log.ComLoggerClient.Error("DescribeCertificateState error: ", err.Error())
			} else {
				updateCertificateState(k, state, cert, key)
			}
		}
	})
	c.Start()
}

func certificateTimeFunc() {
	var certificates []model.Certificate
	if err := model.DdClient.
		Where("auto_renew = ?", 1).
		Where("end_time IS NOT NULL").
		Where("order_id IS NOT NULL").
		Find(&certificates).Error; err != nil {
		log.ComLoggerClient.Error("获取全部已经验证成功的证书且开启了自动续期失败: ", err.Error())
		return
	}

	for i := range certificates {
		v := certificates[i]
		// 判断到期时间减去当前时间是否大于续期时间
		renewalDays := v.RenewTime
		currentTime := time.Now()
		timeRemaining := v.EndTime.Sub(currentTime)
		renewalDuration := time.Duration(renewalDays) * 24 * time.Hour

		if timeRemaining <= renewalDuration {
			log.ComLoggerClient.Info(fmt.Sprintf("证书到期时间小于等于续期时间，开启续期: id = %d, 域名 = %s", v.ID, v.Domain))
			// 判断是否能续期
			if err := util.CanApplyCertificate(); err != nil {
				log.ComLoggerClient.Error(fmt.Sprintf("续期失败: id = %d, 域名 = %s, error = %s", v.ID, v.Domain, err.Error()))
				// 没有额度，直接 return，后面的肯定也申请不了
				return
			}
			if orderId, err := util.CreateCertificateForPackageRequest(v.Domain); err != nil {
				log.ComLoggerClient.Error(fmt.Sprintf("续期失败: id = %d, 域名 = %s, error = %s", v.ID, v.Domain, err.Error()))
				continue
			} else {
				v.OrderId = &orderId
			}
			v.State = "domain_verify"
			v.Validate = false

			if err := model.DdClient.Save(v).Error; err != nil {
				log.ComLoggerClient.Error(fmt.Sprintf("续期失败: id = %d, 域名 = %s, error = %s", v.ID, v.Domain, err.Error()))
				continue
			}

			AddCertificateState(v.ID, *v.OrderId)
		} else {
			log.ComLoggerClient.Info(fmt.Sprintf("证书到期时间大于续期时间，不续期: id = %d, 域名 = %s", v.ID, v.Domain))
		}
	}
}

func initCertificateTimeCron() {
	c := cron.New()
	c.AddFunc("0 0 1 * *", certificateTimeFunc)
	// 初始手动执行下，异步一下
	go certificateTimeFunc()
	c.Start()
}

func updateCertificateState(id uint, state string, cert string, key string) {
	var certificate model.Certificate
	if err := model.DdClient.First(&certificate, id).Error; err != nil {
		log.ComLoggerClient.Error("更新证书状态失败: ", err.Error())
		// 数据被删了需要删除定时任务
		if err == gorm.ErrRecordNotFound {
			RemoveCertificateState(id)
		}
		return
	}

	certificate.State = state

	if state == "certificate" {
		certificate.CertContent = cert
		certificate.KeyContent = key
		// 还需要查询证书到期时间
		s, e, err := util.GetCertificateTime(*certificate.OrderId)
		if err != nil {
			log.ComLoggerClient.Error("获取证书到期时间失败: ", err.Error())
		}
		certificate.StartTime = &s
		certificate.EndTime = &e
		go util.DeployCertificate(cert, key, certificate.CertPath, certificate.KeyPath, certificate.Command, certificate.Domain)
	}

	if state == "domain_verify" && !certificate.Validate {
		// 插入一条 DNS，这里的 cert, key 是表示 主机记录，记录值
		if err := util.CreateDnsRecord(certificate.Domain, cert, key); err != nil {
			log.ComLoggerClient.Error(fmt.Sprintf(
				"创建 DNS 记录失败: id = %d, 域名 = %s, error = %s",
				certificate.ID, certificate.Domain, err.Error()),
			)
		} else {
			log.ComLoggerClient.Info(fmt.Sprintf("创建 DNS 记录成功: id = %d, 域名 = %s", certificate.ID, certificate.Domain))
			certificate.Validate = true
		}
	}

	if err := model.DdClient.Save(&certificate).Error; err != nil {
		log.ComLoggerClient.Error("保存证书状态失败: ", err.Error())
	}

	if state == "certificate" || state == "verify_fail" {
		RemoveCertificateState(id)
	}
	log.ComLoggerClient.Info(fmt.Sprintf("更新证书状态成功: id = %d, state = %s", id, state))
}

func AddCertificateState(id uint, orderId int64) {
	certificateStateMap[id] = orderId
}

func RemoveCertificateState(id uint) {
	delete(certificateStateMap, id)
}
