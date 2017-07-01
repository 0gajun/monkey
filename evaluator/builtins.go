package evaluator

import (
	"unicode/utf8"

	"github.com/0gajun/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if str, ok := args[0].(*object.String); ok {
				length := utf8.RuneCountInString(str.Value)
				return &object.Integer{Value: int64(length)}
			}

			return newError("argument to `len` not supported, got %s", args[0].Type())
		},
	},
}
