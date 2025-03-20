package xerr

// SUCCESS 成功返回
const SUCCESS uint32 = 200
const ErrorCommon uint32 = 1
const LogicSuccess uint32 = 0

/** 全局错误码 **/

const ServerFail uint32 = 500
const RequestParamError uint32 = 400
const UNAUTHORIZED uint32 = 401
const FORBIDDEN uint32 = 403
const RouteNotFound uint32 = 404
const RouteNotMatch uint32 = 405
const PasswordIncorrect uint32 = 406

const DBError uint32 = 600

const (
	ParamMiss          uint32 = 100001
	RecordNotFound     uint32 = 100002
	RepeatRequest      uint32 = 100003
	FileMiss           uint32 = 100004
	SystemError        uint32 = 100005
	SystemBusyError    uint32 = 100006
	RequestError       uint32 = 100007
	InvalidToken       uint32 = 100008
	ParamError         uint32 = 100009
	FileTypeError      uint32 = 100010
	SensitiveError     uint32 = 100011
	MessageLengthError uint32 = 100012
	NotSupportError    uint32 = 100013
)

const (
	ConfigExist          uint32 = 300001
	ConfigOpenAiKeyEmpty uint32 = 300002
	ConfigEmpty          uint32 = 300003
)

const (
	LoginMobileError            uint32 = 200001
	LoginVerifyCodeError        uint32 = 200002
	LoginAccountNotExist        uint32 = 200003
	LoginAccountOrPasswordError uint32 = 200004
	LoginAccountExist           uint32 = 200005
	LoginCaptchaError           uint32 = 200006
	SourceError                 uint32 = 200007
	MobileError                 uint32 = 200008
	LoginInvitationNotExist     uint32 = 200009
	LoginInvitationError        uint32 = 200010
	LoginInvitationIsMyself     uint32 = 200011
	LoginTokenError             uint32 = 200012
	LoginMiss                   uint32 = 200013
	NotificationNotExist        uint32 = 200014
)

const (
	CourseNotFound                      uint32 = 300001
	ChapterNotFound                     uint32 = 300002
	LearningRecordNotFound              uint32 = 300003
	LearningRecordDetailNotFound        uint32 = 300004
	LearningRecordDetailPlayTimeTooLong uint32 = 300005
)

const (
	InvalidSubscriptionType uint32 = 400001
	OrderAlreadyExists      uint32 = 400002
	OrderNotFound           uint32 = 400003
	OrderStatusError        uint32 = 400004
	PaymentFailed           uint32 = 400005
	InvalidOrderType        uint32 = 400006
	OrderUsed               uint32 = 400007
)
