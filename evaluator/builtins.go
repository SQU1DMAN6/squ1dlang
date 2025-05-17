package evaluator

import (
	"fmt"
	"squ1d/object"
	// "unicode/utf8"
)

var builtins = map[string]*object.Builtin{
	"cat": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("Argument to `cat` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"others": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `others` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				// newElements := make([]object.Object, length-1, length-1)
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"add": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `add` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			// newElements := make([]object.Object, length+1, length+1)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"write": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		},
	},
}
