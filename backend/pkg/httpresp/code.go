package httpresp

const (
	Success           Code = "BE0000"
	ServerError       Code = "BE0001"
	BadRequest        Code = "BE0002"
	InvalidRequest    Code = "BE0004"
	Failed            Code = "BE0073"
	Pending           Code = "BE0050"
	InvalidInputParam Code = "BE0032"
	DuplicateUser     Code = "BE0033"
	NotFound          Code = "BE0034"

	Unauthorized   Code = "BE0502"
	Forbidden      Code = "BE0503"
	GatewayTimeout Code = "BE0048"
)

type Code string

var codeMap = map[Code]string{
	Success:           "success",
	Failed:            "failed",
	Pending:           "pending",
	BadRequest:        "bad or invalid request",
	Unauthorized:      "Unauthorized Token",
	GatewayTimeout:    "Gateway Timeout",
	ServerError:       "Internal Server Error",
	InvalidInputParam: "Other invalid argument",
	DuplicateUser:     "duplicate user",
	NotFound:          "Not found",
}

func (c Code) AsString() string {
	return string(c)
}

func (c Code) GetStatus() string {
	switch c {
	case Success:
		return "SUCCESS"

	default:
		return "FAILED"
	}
}

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetVersion() string {
	return "1"
}
