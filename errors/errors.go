package errors

//system
var (
	INTERNAL_ERROR = DataError{
		Code:    -999999,
		Message: "Internal system error, please report to system admin",
	}
)

//common
var (
	REQ_SIZE_TOO_BIG    = Error(10000, "Request body too big")
	ILLEGAL_REQUEST     = Error(10001, "Illegal request")
	REQ_DATA_INVALID    = Error(10002, "Request data invalid")
	API_VERSION_TOO_LOW = Error(10003, "API version too low")
	REQUEST_RATE_LIMIT  = Error(10004, "Request rate limit")
)

//data service
var (
	CLIENT_NOT_FOUND = Error(20000, "Client not found")
)

//gate service
var ()

//notifer service
var ()

//rest_api service
var ()

//session service
var ()
