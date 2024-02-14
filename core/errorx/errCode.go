package errorx

const (
	errCodeBase      = 0x00000000
	errCodeRpcServer = 0x00010000
	errCodeApiServer = 0x00020000
	errCodeRedis     = 0x00030000
	errCodeDB        = 0x00040000
	errCodeMq        = 0x00050000
	errMask          = 0xFFFF
)

const (
	ErrCodeOK       ErrCode = errCodeBase + iota
	ErrCodeStandard         //默认缺少指定error-code类型
	ErrCodeJsonErr
	ErrCodeNotAuthorized
	ErrCodeBusy
	ErrCodeTimeOut
	ErrCodeParamsErr
	ErrCodeNotFound
	ErrCodeNetWorkErr
	ErrCodeUnknown
	ErrCodeLogPathNotSet
	ErrCodeLogFileClosed
	ErrCodeWatcherFileStop
	ErrCodeAssertionValue
	ErrCodeGetPathValue
	ErrCodeMsgTypeMismatch
	ErrCodeBaseMax = errCodeBase | errMask
)

const (
	ErrCodeRedisConnFailed ErrCode = errCodeRedis + iota
	ErrCodeRedisWaitCASTimeout
	ErrCodeRedisKeyIsEmpty
	ErrCodeRedisCacheNotMiss
	ErrCodeRedisCacheDelFailed
)

const (
	ErrCodeDbConnFailed ErrCode = errCodeDB + iota
	ErrCodeDBQueryFailed
	ErrCodeDBQueryCountRepeat
)

const (
	ErrCodeMqRabbitConnFailed ErrCode = errCodeMq + iota
	ErrCodeMqRedisConnFailed
	ErrCodeMqALLConnFailed
)

const (
	ErrRpcCodeParentGroupNotExist ErrCode = errCodeRpcServer + iota
	ErrRpcServerNotFound
	ErrRpcCodeGroupNotEmpty
	ErrRpcCodeConfIsBusy2
	ErrRpcCodeMobileRepeat
	ErrRpcCodeGroupIndepFailed
	ErrRpcCodeGroupLevelTooBig
	ErrRpcCodeCapacityTooBig
	ErrRpcCodeEmailRepeat
	ErrRpcCodeDelTopGroup
	ErrRpcCodeAgentMax = errCodeRpcServer | errMask
)

const (
	ErrApiCodeParentGroupNotExist ErrCode = errCodeApiServer + iota
	ErrApiCodeShouldBindJSON
	ErrApiCodeGroupNotEmpty
	ErrApiCodeConfIsBusy2
	ErrApiCodeMobileRepeat
	ErrApiCodeGroupIndepFailed
	ErrApiCodeGroupLevelTooBig
	ErrApiCodeCapacityTooBig
	ErrApiCodeEmailRepeat
	ErrApiCodeDelTopGroup
	ErrApiCodeAgentMax = errCodeApiServer | errMask
)

type ErrCode uint64

var errMap map[ErrCode]*Status

func init() {
	errMap = map[ErrCode]*Status{
		ErrCodeOK:              {Code: uint64(ErrCodeOK), Message: "Ok"},
		ErrCodeStandard:        {Code: uint64(ErrCodeStandard), Message: "not standard errors"},
		ErrCodeJsonErr:         {Code: uint64(ErrCodeJsonErr), Message: "JsonErr"},
		ErrCodeNotAuthorized:   {Code: uint64(ErrCodeNotAuthorized), Message: "NotAuthorized"},
		ErrCodeUnknown:         {Code: uint64(ErrCodeUnknown), Message: "UnknownMethod"},
		ErrCodeBusy:            {Code: uint64(ErrCodeBusy), Message: "Busy"},
		ErrCodeTimeOut:         {Code: uint64(ErrCodeTimeOut), Message: "TimeOut"},
		ErrCodeParamsErr:       {Code: uint64(ErrCodeParamsErr), Message: "ParamsErr"},
		ErrCodeNotFound:        {Code: uint64(ErrCodeNotFound), Message: "NotFound"},
		ErrCodeNetWorkErr:      {Code: uint64(ErrCodeNetWorkErr), Message: "NetWorkErr"},
		ErrCodeLogPathNotSet:   {Code: uint64(ErrCodeLogPathNotSet), Message: "LogPathNotSet"},
		ErrCodeLogFileClosed:   {Code: uint64(ErrCodeLogFileClosed), Message: "LogFileClosed"},
		ErrCodeWatcherFileStop: {Code: uint64(ErrCodeWatcherFileStop), Message: "WatcherFileStop"},
		ErrCodeAssertionValue:  {Code: uint64(ErrCodeAssertionValue), Message: "type assertion to value failed"},
		ErrCodeGetPathValue:    {Code: uint64(ErrCodeGetPathValue), Message: "config assertion getPath value failed"},
		ErrCodeMsgTypeMismatch: {Code: uint64(ErrCodeMsgTypeMismatch), Message: "eventBus msg TypeMismatch"},

		//redis
		ErrCodeRedisConnFailed:     {Code: uint64(ErrCodeRedisConnFailed), Message: "Redis Conn is failed"},
		ErrCodeRedisWaitCASTimeout: {Code: uint64(ErrCodeRedisWaitCASTimeout), Message: "Redis TryLock wait CAS failed"},
		ErrCodeRedisKeyIsEmpty:     {Code: uint64(ErrCodeRedisKeyIsEmpty), Message: "Redis Key isEmpty"},
		ErrCodeRedisCacheNotMiss:   {Code: uint64(ErrCodeRedisCacheNotMiss), Message: "Redis Cache NotMiss"},
		ErrCodeRedisCacheDelFailed: {Code: uint64(ErrCodeRedisCacheDelFailed), Message: "Redis Cache DelFailed"},

		//db
		ErrCodeDbConnFailed:       {Code: uint64(ErrCodeDbConnFailed), Message: "db conn is Failed"},
		ErrCodeDBQueryFailed:      {Code: uint64(ErrCodeDBQueryFailed), Message: "db Query Failed"},
		ErrCodeDBQueryCountRepeat: {Code: uint64(ErrCodeDBQueryCountRepeat), Message: "db Query Count Repeat"},

		//mq
		ErrCodeMqRabbitConnFailed: {Code: uint64(ErrCodeMqRabbitConnFailed), Message: "Mq Rabbit ConnFailed"},
		ErrCodeMqRedisConnFailed:  {Code: uint64(ErrCodeMqRedisConnFailed), Message: "Mq Redis ConnFailed"},
		ErrCodeMqALLConnFailed:    {Code: uint64(ErrCodeMqALLConnFailed), Message: "Mq ALL ConnFailed"},

		//rpc
		ErrRpcServerNotFound: {Code: uint64(ErrRpcServerNotFound), Message: "ServerNotFound"},

		//api
		ErrApiCodeShouldBindJSON: {Code: uint64(ErrApiCodeShouldBindJSON), Message: "ErrApiCodeShouldBindJSON"},
	}
}

func (e ErrCode) Error() string {
	if v, ok := errMap[e]; ok {
		return v.Message
	}
	return "OK"
}

func (e ErrCode) Code() uint64 {
	return uint64(e)
}

func isErrCode(err error) bool {

	se, ok := err.(ErrCode)
	if se > 0 {
		_, ok = errMap[err.(ErrCode)]
		return ok
	}
	return ok
}

func isStatus(err error, opts ...StatusFunc) (*Status, bool) {

	var (
		isFound bool
		ss      *Status
	)

	if se, ok := err.(ErrCode); ok {
		if ss, isFound = errMap[se]; isFound {
			ss.Metadata = make(map[string]string)
			for _, fn := range opts {
				fn(ss)
			}
		}
	}

	return ss, isFound
}
