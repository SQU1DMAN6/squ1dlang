package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"squ1d/object"
	"strconv"
	"strings"
	"unicode/utf8"
)

var builtins = map[string]*object.Builtin{
	"makedir": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			dirname, ok := args[0].(*object.String)

			if !ok {
				return newError("Argument must be string. Got %s", args[0].Type())
			}

			err := os.Mkdir(dirname.Value, 0777)

			if err != nil {
				return newError("Error making directory: %s", err.Error())
			}

			return nil
		},
	},
	"dirmv": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			source, ok1 := args[0].(*object.String)
			destination, ok2 := args[1].(*object.String)

			if !ok1 || !ok2 {
				return newError("Arguments must be strings. Got %s and %s", args[0].Type(), args[1].Type())
			}

			sourceDir, err := os.Stat(source.Value)
			if err != nil {
				return newError("Failed to stat source directory %s: %s", source.Value, err.Error())
			}
			if !sourceDir.IsDir() {
				return newError("Source is not a directory: %s", source.Value)
			}

			err = os.Rename(source.Value, destination.Value)
			if err != nil {
				return newError("Failed to move directory from %s to %s: %s", source.Value, destination.Value, err.Error())
			}

			return nil
		},
	},
	"filemv": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			source, ok1 := args[0].(*object.String)
			destination, ok2 := args[1].(*object.String)

			if !ok1 || !ok2 {
				return newError("Arguments must be strings. Got %s and %s", args[0].Type(), args[1].Type())
			}

			sourceFile, err := os.Stat(source.Value)
			if err != nil {
				return newError("Failed to stat source file %s: %s", source.Value, err.Error())
			}
			if sourceFile.IsDir() {
				return newError("Source is a directory, not a file: %s", source.Value)
			}

			err = os.Rename(source.Value, destination.Value)
			if err != nil {
				return newError("Failed to move file from %s to %s: %s", source.Value, destination.Value, err.Error())
			}

			return nil
		},
	},
	"dircp": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			src, ok1 := args[0].(*object.String)
			dst, ok2 := args[1].(*object.String)

			if !ok1 || !ok2 {
				return newError("Arguments must be strings. Got %s and %s", args[0].Type(), args[1].Type())
			}

			srcInfo, err := os.Stat(src.Value)
			if err != nil {
				return newError("Failed to stat source directory: %s", err.Error())
			}
			if !srcInfo.IsDir() {
				return newError("Source is not a directory: %s", src)
			}

			// Create the destination directory if it doesn't exist
			if err := os.MkdirAll(dst.Value, srcInfo.Mode()); err != nil {
				return newError("Failed to create destination directory: %s", err.Error())
			}

			entries, err := os.ReadDir(src.Value)
			if err != nil {
				return newError("Failed to read source directory: %s", err.Error())
			}

			for _, entry := range entries {
				srcPath := filepath.Join(src.Value, entry.Name())
				dstPath := filepath.Join(dst.Value, entry.Name())

				if entry.IsDir() {
					if err := copyDir(srcPath, dstPath); err != nil {
						return newError("Failed to recursively copy directories: %s", err.Error())
					}
				} else {
					if err := copyFile(srcPath, dstPath); err != nil {
						return newError("Failed to recursively copy files: %s", err.Error())
					}
				}
			}
			return nil
		},
	},
	"filecp": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			source, ok1 := args[0].(*object.String)
			stat, ok2 := args[1].(*object.String)

			if !ok1 || !ok2 {
				return newError("Arguments must be strings. Got %s and %s", args[0].Type(), args[1].Type())
			}

			src, err := os.Open(source.Value)
			if err != nil {
				return newError("Error opening source file: %s", err.Error())
			}
			defer src.Close()

			dst, err := os.Create(stat.Value)
			if err != nil {
				return newError("Error creating destination file: %s", err.Error())
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return newError("Error copying file: %s", err.Error())
			}

			return nil
		},
	},
	"ls": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			input, ok := args[0].(*object.String)

			if !ok {
				return newError("Argument must be string. Got %s", args[0].Type())
			}

			dirPath := input.Value
			entries, err := os.ReadDir(dirPath)

			if err != nil {
				return newError("Error reading directory: %s", err.Error())
			}

			var result []object.Object

			for _, entry := range entries {
				result = append(result, &object.String{Value: entry.Name()})
			}

			return &object.Array{Elements: result}
		},
	},
	"writefile": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2", len(args))
			}

			filePath, ok1 := args[0].(*object.String)
			fileContents, ok2 := args[1].(*object.String)

			if !ok1 || !ok2 {
				return newError("Arguments must be string. Got %s and %s", args[0].Type(), args[1].Type())
			}

			data := []byte(fileContents.Value)
			err := os.WriteFile(filePath.Value, data, 0777)
			if err != nil {
				return newError("Error writing file: %s", err.Error())
			}

			return nil
		},
	},
	"readfile": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			filePath, ok := args[0].(*object.String)

			if !ok {
				return newError("Arguments must be string. Got %s", args[0].Type())
			}

			contents, err := os.ReadFile(filePath.Value)
			if err != nil {
				return newError("Failed to read file: %s", err.Error())
			}

			return &object.String{Value: string(contents)}
		},
	},
	"read": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			prompt, ok := args[0].(*object.String)

			if !ok {
				return newError("Arguments must be string. Got %s", args[0].Type())
			}

			fmt.Print(prompt.Value)

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return newError("Failed to read input: %s", err.Error())
			}

			input = strings.TrimSpace(input)

			var value object.Object
			if intVal, err := strconv.ParseInt(input, 10, 64); err == nil {
				value = &object.Integer{Value: intVal}
			} else {
				value = &object.String{Value: input}
			}

			return value
		},
	},
	"intstr": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments. Got %d, expected 1", len(args))
			}

			intObj, ok := args[0].(*object.Integer)
			if !ok {
				return newError("Argument must be an integer. Got %s", args[0].Type())
			}

			strVal := strconv.Itoa(int(intObj.Value))

			return &object.String{Value: strVal}
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
				return arr.Elements[length-1] //convert from [] to {}
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
	"arraycontains": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong number of arguments. Got %d, expected 2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `contains` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			element := args[1]
			length := len(arr.Elements)

			for i := 0; i < length; i++ {
				fmt.Printf("arraylement %v, element %v", arr.Elements[i], element)
				if arr.Elements[i] == element {
					return &object.Boolean{Value: true}
				}
			}
			return &object.Boolean{Value: false}

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

func copyDir(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcDir, err)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents from %s to %s: %w", src, dst, err)
	}

	return nil
}
