package engine

type Request struct {
	Url    string
	Parser Parser
}

type Parser interface {
	Parse(contents []byte, url string) ParserResult
	Serialize() (name string, args interface{})
}

type ParserFun func(content []byte, url string) ParserResult

type FunParser struct {
	parser ParserFun
	name   string
}

func (f *FunParser) Parse(content []byte, url string) ParserResult {
	return f.parser(content, url)
}

func (f *FunParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}

// 工厂模式，创建函数解析器
func NewFuncParser(fun ParserFun, funcName string) *FunParser {
	return &FunParser{
		parser: fun,
		name:   funcName,
	}
}

type ParserResult struct {
	Items    []Item
	Requests []Request
}

type Item struct {
	Id      string
	Url     string
	PayLoad interface{}
}

func NilParser([]byte) ParserResult {
	return ParserResult{}
}
