package statuscode

const (
	StatusCommonBackendError     = 1000 // Backend error code start from 1000 to 9999
	StatusBindingInputJsonFailed = 1001
	StatusCreateItemFailed       = 1002
	StatusUpdateItemFailed       = 1003
	StatusDeleteItemFailed       = 1004
	StatusSearchItemFailed       = 1005
	StatusItemNotFound           = 1006
)
