package java

//一整个类组装生成

type Field struct {
	Name    string
	Type    string
	Static  bool
	Imports Imports `toml:"imports"`
}

// 静态代码块
type StaticCode struct{}

type Method struct {
	Name      string            `toml:"name"`
	Static    bool              `toml:"static"`
	Code      string            `toml:"code"`
	Imports   Imports           `toml:"imports"`
	Variables map[string]string `toml:"variables"`
}

type Imports struct {
	Classes []string `toml:"classes"`
}

type Config struct {
	Methods []Method `toml:"methods"`
}

type Class struct {
	Methods    []Method
	Fields     []Field
	Imports    Imports
	ClassName  string
	StaticCode StaticCode
	Extends    string
	Implements string
}
