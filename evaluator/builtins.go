package evaluator

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"squ1d/object"
	"strconv"
	"strings"
	"unicode/utf8"
)

// rand3(1, 2, 3)
var builtins = map[string]*object.Builtin{
	"read": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			prompt, ok := args[0].(*object.String)

			if !ok {
				return newError("Arguments must be strings. Got %s and %s", args[0].Type(), args[1].Type())
			}

			fmt.Print(prompt.Value)

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return newError("Failed to read input: %s", err.Error())
			}

			input = strings.TrimSpace(input)

			// Try to parse as integer
			var value object.Object
			if intVal, err := strconv.ParseInt(input, 10, 64); err == nil {
				value = &object.Integer{Value: intVal}
			} else {
				value = &object.String{Value: input}
			}

			return value
		},
	},
	"tpint": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			strObj, ok := args[0].(*object.String)
			if !ok {
				return newError("Argument must be a string. Got %s", args[0].Type())
			}

			intVal, err := strconv.ParseInt(strObj.Value, 10, 64)
			if err != nil {
				return newError("Failed to convert to integer: %s", err.Error())
			}

			return &object.Integer{Value: intVal}
		},
	},
	"rand": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
	"sepr": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			strObj, ok1 := args[0].(*object.String)
			indexObj, ok2 := args[1].(*object.Integer)

			if !ok1 || !ok2 {
				return newError("Arguments must be (string, integer). Got %s and %s", args[0].Type(), args[1].Type())
			}

			parts := strings.Fields(strObj.Value)
			idx := int(indexObj.Value)

			if idx < 0 || idx >= len(parts) {
				return newError("Index out of bounds. Got %d, but only %d parts", idx, len(parts))
			}

			return &object.String{Value: parts[idx]}
		},
	},
	"tp": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
	"add": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
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
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}

			fmt.Println()
			return nil
		},
	},
}
