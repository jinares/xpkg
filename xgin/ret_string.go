package xgin

import "google.golang.org/grpc/codes"

type (
	RetString struct {
		data string `json:"data"`
		ct   ContentType
	}
)

func (c *RetString) GetCode() codes.Code {
	return codes.OK
}
func (c *RetString) GetRet() (string, ContentType) {
	return c.data, c.ct
}

func (c *RetString) SetTrace(str string) IRet {
	//c.Trace = str
	return c
}
func NewRetString(data string) *RetString {
	return NewRetStringv2(data, TEXT_HTML)
}

func NewRetStringv2(data string, ct ContentType) *RetString {
	return &RetString{
		data: data,
		ct:   ct,
	}
}

/*
func (t *RetString) UnmarshalJSON(b []byte) error

	return nil
}

func (t RetString) MarshalJSON() ([]byte, error) {

	return []byte(""), nil
}
*/
