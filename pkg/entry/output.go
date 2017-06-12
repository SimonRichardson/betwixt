package entry

// Document defines a struct of all the possible things that a http request
// could encounter.
type Document struct {
	Method      *String
	Status      *Status
	URL         *URL
	Params      *Map
	ReqHeaders  *Map
	ReqBody     *String
	RespHeaders *Map
	RespBody    *String
}
