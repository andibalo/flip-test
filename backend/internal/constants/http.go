package constants

const (
	RequestID      string = `X-Request-Id`
	RequestMethod  string = `x-request-method`
	RequestScheme  string = `x-request-scheme`
	KeyServerRoute string = `x-key-server-route`
	ForwardedFor   string = `x-forwarded-for`
	XClientID      string = `X-Client-Id`
	XClientVersion string = `X-Client-Version`
	XUserEmail     string = `x-user-email`

	AppToken string = `x-app-token`

	LangEN string = `EN`
	LangID string = `ID`

	UserAgent             string = `User-Agent`
	ContentAccept         string = `Accept`
	ContentType           string = `Content-Type`
	ContentJSON           string = `application/json`
	ContentXML            string = `application/xml`
	ContentFormURLEncoded string = `application/x-www-form-urlencoded`

	CacheControl        string = `Cache-Control`
	CacheNoCache        string = `no-cache`
	CacheNoStore        string = `no-store`
	CacheMustRevalidate string = `must-revalidate`
)
