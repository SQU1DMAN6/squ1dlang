package evaluator

import (
	"fmt"
	"squ1d/object"
	"unicode/utf8"
	"math/rand"
)

var builtins = map[string]*object.Builtin{
	"rand": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			min, ok1 := args[0].(*object.Integer)
			max, ok2 := args[1].(*object.Integer)

			if !ok1 || !ok2 {
				return newError("Wrong argument type, expected integer, got %T", args[0].Type())
			}

			if min.Value >= max.Value {
				return newError("First argument must be less than second argument")
			}

			rangeInt := int(max.Value - min.Value)

			randNum := rand.Intn(rangeInt) + int(min.Value)

			return &object.Integer{Value: int64(randNum)}
		},
	},
	"tp": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				return &object.String{Value: "array"}
			case *object.String:
				return &object.String{Value: "string"}
			case *object.Hash:
				return &object.String{Value: "hash"}
			case *object.Integer:
				return &object.String{Value: "integer"}
			case *object.Boolean:
				return &object.String{Value: "boolean"}
			case *object.Function:
				return &object.String{Value: "function"}
			default:
				return &object.String{Value: "null"}
			}
		},
	},
	"cat": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
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
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"write": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}

			fmt.Println()
			return nil
		},
	},
}