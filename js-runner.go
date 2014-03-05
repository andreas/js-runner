package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/idada/v8.go"
	"os"
	"reflect"
)

type Input struct {
	Script string `json:"script"`
	Variables map[string]interface{} `json:"variables"`
}

type Output struct {
	Type string `json:"type"`
	Value *json.RawMessage `json:"value"`
}

type Exception struct {
	Type string `json:"type"`
	Description string `json:"description"`
}

func typeName(v *v8.Value) string {
	switch {
	case v.IsUndefined():
		return "undefined"
	case v.IsString():
		return "string"
	case v.IsBoolean():
		return "boolean"
	case v.IsNumber():
		return "number"
	case v.IsFunction():
		return "function"
	case v.IsDate():
		return "date"
	}
	
	return "object"
}

func main() {
	var input Input
	
	// Read and parse input specification from stdin
	bytes, err := ioutil.ReadAll(os.Stdin); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	err = json.Unmarshal(bytes, &input); if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}

	// Setup V8
	engine := v8.NewEngine()
	context := engine.NewContext(nil)
	
	context.Scope(func(cs v8.ContextScope) {
		exception := cs.TryCatch(true, func() {
			// Prepare global context
			global := cs.Global()
			
			for k, v := range input.Variables {
				global.SetProperty(k, engine.GoValueToJsValue(reflect.ValueOf(v)), v8.PA_None)
			}

			// Run script			
			result := cs.Eval(input.Script)
			if result != nil {
				// Output result
				resultJSON := json.RawMessage(v8.ToJSON(result))
				output := Output{typeName(result), &resultJSON}
				outputJSON, err := json.Marshal(output); if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Stdout.Write(outputJSON)
			}
		})
		
		// Something went wrong!
		if exception != "" {
			output, err := json.Marshal(Exception{"RuntimeException", exception}); if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Stdout.Write(output)
		}
	})
}
