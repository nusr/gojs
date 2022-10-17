package token

import (
	"fmt"
	"strconv"
)

type Type int

const (
	LeftParen       Type = iota // (
	RightParen                  // )
	LeftBrace                   // {
	RightBrace                  // }
	LeftSquare                  // [
	RightSquare                 // ]
	Comma                       // ,
	Dot                         // .
	Minus                       // -
	MinusEqual                  // -=
	MinusMinus                  // --i
	Plus                        // +
	PlusEqual                   // +=
	PlusPlus                    // ++
	Semicolon                   // ;
	Colon                       // :
	Slash                       // /
	SlashEqual                  // /=
	Star                        // *
	StarEqual                   // *=
	StarStar                    // **
	StarStarEqual               // **=
	Percent                     // %
	PercentEqual                // %=
	Mark                        // ?
	Bang                        // one or two character tokens !
	BangEqual                   // !=
	BangEqualEqual              // !==
	Equal                       // =
	EqualEqual                  // ==
	EqualEqualEqual             // ===
	Greater                     // >
	GreaterEqual                // >=
	Less                        // <
	LessEqual                   // <=
	Identifier                  // Literals
	String
	Float64
	Int64
	And      // keywords
	AndEqual // &&=
	Class
	Else
	False
	True
	Function
	For
	If
	Null // null
	Or
	OrEqual                    // ||=
	BitAnd                     // &
	BitXOr                     // ^
	BitXOrEqual                // ^=
	BitAndEqual                // &=
	BitOr                      // |
	BitOrEqual                 // |=
	BitNot                     // ~
	BitLeftShift               // <<
	BitLeftShiftEqual          // <<=
	BitRightShift              // >>
	BitRightShiftEqual         // >>=
	BitUnsignedRightShift      // >>>
	BitUnsignedRightShiftEqual // >>>=
	Return
	Super
	This
	Static // static
	Var    // variable
	Do     // do
	While  // while
	New    // new
	EOF    // end
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
