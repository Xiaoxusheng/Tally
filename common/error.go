package common

type Errors struct {
	Code     int
	ErrorMsg string
}

func (e *Errors) Error() string {
	return e.ErrorMsg
}

func (e *Errors) Replace(msg string) {
	e.ErrorMsg = msg
}

func New(code int, msg string) Errors {
	return Errors{
		Code:     code,
		ErrorMsg: msg,
	}
}
