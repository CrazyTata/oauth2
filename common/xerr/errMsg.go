package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[SUCCESS] = "成功"
	message[ServerFail] = "调用失败请稍后重试"
	message[RequestParamError] = "参数错误"
	message[UNAUTHORIZED] = "无效的token"
	message[FORBIDDEN] = "权限不足"
	message[PasswordIncorrect] = "密码错误"
	message[RouteNotFound] = "请求资源不存在"
	message[RouteNotMatch] = "请求方式错误"
	message[DBError] = "数据库繁忙,请稍后再试"

	//common
	message[RecordNotFound] = "record not found"
	message[RepeatRequest] = "please do not repeat the request"
	message[ParamMiss] = "param miss"
	message[FileMiss] = "no such file"
	message[FileTypeError] = "file type error"
	message[SystemError] = "系统错误"
	message[SystemBusyError] = "系统繁忙，请稍后再试～"
	message[RequestError] = "请求参数错误，请稍后再试～"
	message[InvalidToken] = "无效的token"
	message[ParamError] = "param error"
	message[SensitiveError] = "您的消息中含有敏感词信息，请重新输入"
	message[MessageLengthError] = "您的消息太长，请重新输入"
	message[NotSupportError] = "暂不支持此操作，请联系客服处理"

	//config
	message[ConfigExist] = "当前配置在系统中已存在"
	message[ConfigOpenAiKeyEmpty] = "未配置key"
	message[ConfigEmpty] = "缺少配置"

	//login
	message[LoginMobileError] = "mobile error"
	message[LoginTokenError] = "手机号一键登录失败:token error"
	message[LoginVerifyCodeError] = "verify code error"
	message[LoginAccountNotExist] = "account dose not exist"
	message[LoginAccountOrPasswordError] = "account or password error"
	message[LoginAccountExist] = "account is exist"
	message[LoginCaptchaError] = "captcha error"
	message[SourceError] = "来源类型错误"
	message[MobileError] = "手机号错误"
	message[LoginInvitationNotExist] = "邀请码不存在"
	message[LoginInvitationError] = "邀请码错误"
	message[LoginInvitationIsMyself] = "不能使用自己的邀请码"
	message[LoginMiss] = "缺少登录信息"
	message[NotificationNotExist] = "通知不存在"
	//learning
	message[CourseNotFound] = "课程不存在"
	message[ChapterNotFound] = "章节不存在"
	message[LearningRecordNotFound] = "学习记录不存在"
	message[LearningRecordDetailNotFound] = "章节学习记录不存在"
	message[LearningRecordDetailPlayTimeTooLong] = "章节学习记录播放时间过长"

	//order
	message[InvalidSubscriptionType] = "无效的订阅类型"
	message[OrderAlreadyExists] = "订单已存在"
	message[OrderNotFound] = "订单不存在"
	message[OrderStatusError] = "订单状态错误"
	message[PaymentFailed] = "支付失败"
	message[InvalidOrderType] = "无效的订单类型"
	message[OrderUsed] = "订单已使用"
}

func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
