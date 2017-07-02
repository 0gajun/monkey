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

			switch arg := args[0].(type) {
			case *object.String:
				length := utf8.RuneCountInString(arg.Value)
				return &object.Integer{Value: int64(length)}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if array, ok := args[0].(*object.Array); ok {
				if len(array.Elements) == 0 {
					return NULL
				}
				return array.Elements[0]
			}

			return newError("argument to `first` not supported, got %s", args[0].Type())
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if array, ok := args[0].(*object.Array); ok {
				length := len(array.Elements)
				if length == 0 {
					return NULL
				}
				return array.Elements[length-1]
			}

			return newError("argument to `last` not supported, got %s", args[0].Type())
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch array := args[0].(type) {
			case *object.Array:
				length := len(array.Elements)
				if length == 0 {
					return NULL
				}

				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, array.Elements[1:length])
				return &object.Array{Elements: newElements}

			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}

		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			length := len(array.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, array.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
}
