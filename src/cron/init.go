package cron

func Init() {
	initCertificateStateCron()
	initCertificateTimeCron()
	initGetUserCertificateOrder()
}
