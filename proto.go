package main
import (
	"os"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	files = kingpin.Arg("proto", "Protbuf files.").Required().Strings()
)

type Constant struct {
	Pos       lexer.Position
	Str       *string   `  @String`
	Float     *float64  `| @Float`
	Int       *int64    `| @Int`
	Bool      *bool     `| ( @"true" | "false" )`
	Reference *string   `| @Ident { @"." @Ident }`
	Minus     *Constant `| "-" @@`
}

type Import struct {
	Pos  lexer.Position
	Kind string `"import" @["weak" | "public"]`
	Name string `@String ";"`
}

type Package struct {
	Pos  lexer.Position
	Name string `"package" @Ident {@"." @Ident} ";"`
}

type Option struct {
	Pos      lexer.Position
	Name     string    `"option" (@Ident | "(" @Ident {@"." @Ident} ")") { "." @Ident }`
	Constant *Constant `"=" @@ ";"`
}

type ValueOption struct {
	Pos      lexer.Position
	Name     string    `(@Ident | "(" @Ident {@"." @Ident} ")") { "." @Ident }`
	Constant *Constant `"=" @@ `
}

type EnumField struct {
	Pos     lexer.Position
	Name    string         `@Ident "="`
	Value   int            `@Int`
	Options []*ValueOption `[ "[" @@ { "," @@ } "]" ] ";"`
}

type Enum struct {
	Pos     lexer.Position
	Name    string       `"enum" @Ident "{"`
	Options []*Option    `{ @@  `
	Cases   []*EnumField `| @@ | ";" } "}"`
}

type Field struct {
	Pos          lexer.Position
	Repeated     bool           `(@"repeated"`
	Type         string         `(@Ident {@"." @Ident}) | (@Ident {@"." @Ident}))`
	Name         string         `@Ident`
	FieldNumber  int            `"=" @Int`
	FieldOptions []*ValueOption `[ "[" @@ { "," @@} "]" ] ";"`
}

type OneofField struct {
	Pos          lexer.Position
	Type         string         `@Ident {@"." @Ident}`
	Name         string         `@Ident`
	FieldNumber  int            `"=" @Int`
	FieldOptions []*ValueOption `[ "[" @@ { "," @@} "]" ] ";"`
}

type Oneof struct {
	Pos    lexer.Position
	Name   string        `"oneof" @Ident "{"`
	Fields []*OneofField ` { @@ } "}"`
}

type MapField struct {
	Pos          lexer.Position
	KeyType      string         `"map" "<" @("int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string") ","`
	ValueType    string         `@Ident {@"." @Ident} ">"`
	Name         string         `@Ident`
	FieldNumber  int            `"=" @Int`
	FieldOptions []*ValueOption `[ "[" @@ { "," @@} "]" ] ";"`
}

type Range struct {
	Pos   lexer.Position
	From  int  `@Int`
	To    int  `["to" (@Int | `
	ToMax bool `        @"max")]`
}

type Reserved struct { // not completley correct as it can handle both now
	Pos      lexer.Position
	reserved string   `"reserved"`
	Fields   []string `[@String {"," @String}]`
	Ranges   []*Range `[@@ { "," @@ }] ";"`
}

type Message struct {
	Pos       lexer.Position
	Name      string      `"message" @Ident "{"`
	Enums     []*Enum     `{ @@ `
	Messages  []*Message  ` | @@ `
	MapFields []*MapField ` | @@`
	Options   []*Option   ` | @@`
	Oneofs    []*Oneof    ` | @@ `
	Reserved  []*Reserved ` | @@`
	Fields    []*Field    ` | @@ } "}"`
}

type Rpc struct {
	Pos           lexer.Position
	Name          string    `"rpc" @Ident "("`
	RequestStream bool      `[@"stream"]`
	RequestType   string    `@(["."] Ident { "." Ident }) ")"`
	ReturnStream  bool      `"returns" "(" [@"stream"]`
	ReturnType    string    `@(["."] Ident { "." Ident }) ")"`
	Options       []*Option `(("{" { @@ } "}") | ";")`
}

type Service struct {
	Pos     lexer.Position
	Name    string    `"service" @Ident "{"`
	Options []*Option `{ @@ `
	Rpcs    []*Rpc    ` | @@ } "}"`
}

type Proto struct {
	Pos      lexer.Position
	Syntax   string     `"syntax" "=" @String ";"`
	Imports  []*Import  `{ @@ `
	Packages []*Package `  | @@ `
	Options  []*Option  `  | @@ `
	Enums    []*Enum    `  | @@ `
	Services []*Service `  | @@ `
	Messages []*Message `  | @@ }`
}

func main() {
	kingpin.Parse()

	parser, err := participle.Build(&Proto{}, nil)
	println(parser.String())
	kingpin.FatalIfError(err, "")

	for _, file := range *files {
		proto := &Proto{}
		r, err := os.Open(file)
		kingpin.FatalIfError(err, "")
		err = parser.Parse(r, proto)
		kingpin.FatalIfError(err, "")
		spew.Dump(proto)

	}
}
