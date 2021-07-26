package engine

type Request struct {
	Url       string
	ParserFun func([]byte) ParserResult
}

type ParserResult struct {
	Items    []Item
	Requests []Request
}

type Item struct {
	Id string
	Url string
	PayLoad interface{}
}

func NilParser([] byte) ParserResult {
	return ParserResult{}
}
