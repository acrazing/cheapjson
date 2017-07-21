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
  
  // And you can dump a value to raw struct
  data := value.Value()
  // and this could be json marshal
  _, _ = json.Marshal(data)
}
```

## Benchmark

See <./parser_test.go>

Compare with [go-simplejson](github.com/bitly/go-simplejson):

```text
$ go test -v -bench=. ./parser_test.go

=== RUN   TestUnmarshal
--- PASS: TestUnmarshal (0.00s)
=== RUN   TestSimpleJson
--- PASS: TestSimpleJson (0.00s)
=== RUN   TestUnmarshal2
--- PASS: TestUnmarshal2 (0.00s)
BenchmarkUnmarshal-4    	 2000000	       574 ns/op
BenchmarkSimpleJson-4   	 2000000	       801 ns/op
PASS
ok  	command-line-arguments	4.494s
```

## License

MIT


## TODO

- [ ] more unit test
