# [Cheap JSON](https://godoc.org/github.com/acrazing/cheapjson) &middot; [![GoDoc](https://godoc.org/github.com/acrazing/cheapjson?status.svg)](https://godoc.org/github.com/acrazing/cheapjson) [![Build Status](https://travis-ci.org/acrazing/cheapjson.svg?branch=master)](https://travis-ci.org/acrazing/cheapjson)

A arbitrary JSON parser for golang.

- **Standalone**: implement the parser independently for `ECMA-404 The JSON Data Interchange Standard`.
- **Fast**: about two times faster than the package `go-simplejson` which use native `encoding/json` library.
- **Lightweight**: only about 500 rows code for parser include UTF-16 pairs covert to UTF-8 bytes.

## Install

```bash
go get github.com/acrazing/cheapjson
```

## Usage Example

```go
package main

import (
  "github.com/acrazing/cheapjson"
  "encoding/json"
)

func main()  {
  // Unmarshal a bytes slice
  value, err := cheapjson.Unmarshal([]byte("{\"hello\":\"world\", \"int\":12345}"))
  if err != nil {
    panic(err)
  }
  // type check
  if !value.IsObject() {
    panic("parse error")
  }
  // get a child field
  str := value.Get("hello")
  
  // get as string
  println(str.String()) // world
  
  // get as int
  println(value.Get("int").Int()) // 12345
  
  // And any else you can do:
  _ = value.Float() // returns float64
  _ = value.Array() // returns []*Value
  _ = value.Object() // returns map[string]*Value
  
  // WARNING: any of the upon value extract operate
  // need to check the type at first as follow:
  if value.IsObject() {
    // value is a object, and then you can operate:
    _ = value.Object()
  }
  // And there are more type checks
  _ = value.IsObject()
  _ = value.IsArray()
  _ = value.IsNumber() // if is float or int, returns true
  _ = value.IsInt() // just check is int
  _ = value.IsTrue()
  _ = value.IsFalse()
  _ = value.IsBool()
  _ = value.IsNull()
  _ = value.IsString()
  
  // And you can manipulate a value
  value = cheapjson.NewValue()
  value.AsObject(nil) // set as a object
  _ = value.AddField("hello") // if a value is a object, you can call this, else will panic
  value.AsArray(nil) // set as a array
  elem := value.AddElement() // if a value is a array, yu can call this, else will panic
  elem.AsInt(12) // as a int
  elem.AsFloat(232)
  elem.AsBool(true)
  elem.AsNull()
  
  // And you can get a deep path by:
  field := elem.Get("hello", "world", "deep", "3")
  _ = field.Value()
  // Or set a deep path
  // The different between Get and Ensure is that the Get
  // just returns the exists field, if the path does not exist
  // will return nil, and it will covert the path to integer
  // if the node is an array, and the Ensure will force the
  // path to be an object, and if the target path does not exist
  // will auto generate it as a empty node.
  value.Ensure("hello", "world", "deep", "3").AsInt(3)
  
  // And you can dump a value to raw struct
  data := value.Value()
  // and this could be json marshal
  _, _ = json.Marshal(data)
}
```

## Benchmark

See [parser_test.go](./parser_test.go), compare with [go-simplejson](https://github.com/bitly/go-simplejson), which
use the native `encoding/json` library to unmarshal a json. The result is:

- NormalInput(small): about 1.6 times faster
- BigInput: about 4.4 times faster
- DeepInput: about 7 times faster

```bash
go test -bench=. -v ./parser_test.go

# 2017/07/22 12:48:45 big input size: 92338772, normal input size: 763, deep input size: 33976002
# === RUN   TestUnmarshal
# --- PASS: TestUnmarshal (0.00s)
# === RUN   TestSimpleJson
# --- PASS: TestSimpleJson (1.62s)
# BenchmarkUnmarshalBigInput-4                   5         358595392 ns/op
# BenchmarkSimpleJsonBigInput-4                  1        1560047078 ns/op
# BenchmarkUnmarshalNormalInput-4           200000              5372 ns/op
# BenchmarkSimpleJsonNormalInput-4          200000              8593 ns/op
# BenchmarkUnmarshalDeepInput-4                 30          42870590 ns/op
# BenchmarkSimpleJsonDeepInput-4                 5         305351224 ns/op
# PASS
# ok      command-line-arguments  18.314s

```

## License

MIT


## TODO

- [x] more unit test.
- [ ] test the performance about make buffer before handle a string, (will walk the string twice).
