package cheapjson_test

import (
	"testing"
	"."
	"github.com/stretchr/testify/assert"
	"math"
)

func TestNewValue(t *testing.T) {
	object := cheapjson.NewValue()
	object.AsObject(nil)
	assert.Equal(t, true, object.IsObject())
	assert.Equal(t, false, object.IsNull())
	hello := object.AddField("hello")
	hello.AsString("world")
	assert.Equal(t, "world", hello.String())
	assert.Equal(t, "world", object.Get("hello").String())
	integer := object.AddField("integer")
	integer.AsInt(math.MaxInt64)
	assert.Equal(t, int64(math.MaxInt64), integer.Int())
	assert.Equal(t, true, integer.IsInt())
	assert.Equal(t, true, integer.IsNumber())
	double := object.AddField("float")
	double.AsFloat(math.SmallestNonzeroFloat64)
	assert.Equal(t, math.SmallestNonzeroFloat64, double.Float())
	assert.Equal(t, true, double.IsNumber())
	assert.Equal(t, false, double.IsInt())
	assert.Equal(t, false, double.IsString())
	null := object.AddField("null")
	null.AsNull()
	assert.Equal(t, nil, null.Value())
	assert.Equal(t, true, null.IsNull())
}
