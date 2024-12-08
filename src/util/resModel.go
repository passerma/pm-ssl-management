package util

// 错误信息
var errInfo = map[int]string{
	0: "成功",
	1: "错误",
	2: "参数错误",
	3: "数据库错误",
	4: "用户名或密码错误",
	5: "无权限",
}

// 返回信息结构体
type ResModel struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// {功能} 返回正确信息结构体
// {参数} data 数据, info 消息
// {返回} 结构体
func SendSusModel(data interface{}, info ...string) ResModel {
	var res ResModel
	if info != nil {
		res.Msg = info[0]
	} else {
		res.Msg = errInfo[0]
	}
	res.Code = 0
	res.Data = data
	return res
}

// {功能} 返回错误信息结构体
// {参数} data 数据, info 消息
// {返回} 结构体
func SendErrModel(code int, data ...interface{}) ResModel {
	var res ResModel
	if msg, ok := errInfo[code]; ok {
		res.Msg = msg
		res.Code = code
	} else {
		res.Msg = "未知错误"
		res.Code = code
	}
	if data != nil {
		res.Data = data[0]
	} else {
		res.Data = nil
	}
	return res
}
