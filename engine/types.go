package engine

type Request struct {
	Url       string
	ParserFun func([]byte) ParserResult
}

type ParserResult struct {
	Items    []interface{}
	Requests []Request
}

func NilParser([] byte) ParserResult {
	return ParserResult{}
}
