package token

import (
	"fmt"
	"strconv"
)

type Type int

const (
	LeftParen    Type = iota // (
	RightParen               // )
	LeftBrace                // {
	RightBrace               // }
	LeftSquare               // [
	RightSquare              // ]
	Comma                    // ,
	Dot                      // .
	Minus                    // -
	MinusMinus               // --i
	Plus                     // +
	PlusPlus                 // ++
	Semicolon                // ;
	Colon                    // :
	Slash                    // /
	Star                     // *
	Percent                  // %
	Mark                     // ?
	Bang                     // one or two character tokens !
	BangEqual                // !=
	Equal                    // =
	EqualEqual               // ==
	Greater                  // >
	GreaterEqual             // >=
	Less                     // <
	LessEqual                // <=
	Identifier               // Literals
	String
	Float64
	Int64
	And // keywords
	Class
	Else
	False
	True
	Function
	For
	If
	Null // null
	Or
	BitAnd
	BitOr
	Return
	Super
	This
	Var   // variable
	Do    // do
	While // while
	New   // new
	EOF   // end
)

func ConvertAnyToString(text any) string {
	switch data := text.(type) {
	case nil:
		return "null"
	case string:
		return data
	case float64:
		return strconv.FormatFloat(data, 'f', 10, 64)
	case int64:
		return strconv.FormatInt(data, 10)
	case bool:
		{
			if data {
				return "true"
			}
			return "false"
		}
	case fmt.Stringer:
		return data.String()
	default:
		return ""
	}
}
