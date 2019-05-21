package klinreq

type ReqQ struct {
	Url      *string
	Headers  map[string]string
	Method   *string
	NoVerify bool
	Json     *interface{}
}
type ReqBuilder struct {
	ReqQ
}
