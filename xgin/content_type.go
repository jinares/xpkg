package xgin

type (
	ContentType string
)

const (
	APPLICATION_JSON ContentType = "application/json"
	//text/html;text/plain; charset=utf-8
	TEXT_HTML  ContentType = "text/html"
	TEXT_PLAIN ContentType = "text/plain; charset=utf-8"
	//application/javascript
	APPLICATION_JAVASCRIPT ContentType = "application/javascript"
	//application/xml
	APPLICATION_XML ContentType = "application/xml"
	//application/x-www-form-urlencoded
	APPLICATION_URLENCODED ContentType = "application/x-www-form-urlencoded"
)
