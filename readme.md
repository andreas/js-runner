Either download the binary for OSX:

https://github.com/andreas/js-runner/releases/download/untagged-37b72b867e3a7ce8bf2a/js-runner

Or build it yourself:

1. Install `idada/v8.go`: https://github.com/idada/v8.go/
2. `go build main.go`

Then run it:

```
  > echo "{\"script\": \"x+1\", \"variables\": { \"x\": 2 }}" | ./js-runner
  {"type":"number","value":3}
```
