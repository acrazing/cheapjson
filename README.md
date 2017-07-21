# Cheap JSON

A fast, lightweight, struct-less JSON parser for golang.

## Install

```bash
go get github.com/acrazing/cheapjson
```

## Usage

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
  value.AsObject() // set as a object
  _ = value.AddField("hello") // if a value is a object, you can call this, else will panic
  value.AsArray() // set as a array
  elem := value.AddElement() // if a value is a array, yu can call this, else will panic
  elem.AsInt(12) // as a int
  elem.AsFloat(232)
  elem.AsBool(true)
  elem.AsNull()
  
  // And you can get a deep path by:
  field := elem.Get("hello", "world", "deep", "3")
  _ = field.Value()
  
  // Or set a deep path
  value.Ensure("hello", "world", "deep", "3").AsInt(3)
  
  // And you can dump a value to raw struct
  data := value.Value()
  // and this could be json marshal
  _, _ = json.Marshal(data)
}
```

## Benchmark

See [parser_test.go](./parser_test.go), compare with [go-simplejson](https://github.com/bitly/go-simplejson), which
use the native `encoding/json` library to unmarshal a json. The result is half of the time to cost!

```text
$ go test -v -bench=. ./parser_test.go -run NONE

2017/07/22 00:25:39 big input size: 78122
BenchmarkUnmarshal-4                3000            537921 ns/op
BenchmarkSimpleJson-4               2000           1117259 ns/op
PASS
ok      command-line-arguments  4.047s
```

## License

MIT


## TODO

- [ ] more unit test
