package klinreq

type reqQ struct {
	Url     *string
	Headers map[string]string
	Method  *string
	Json    *interface{}
}
type reqBuilder struct {
	reqQ
}
