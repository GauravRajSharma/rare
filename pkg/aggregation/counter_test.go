package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAddition(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	items := val.Items()

	assert.Equal(t, 1, len(items), "Expected length 1")
	assert.Equal(t, "test", items[0].Name)
	assert.Equal(t, int64(1), items[0].Item.count)
}

func collectChan(c chan interface{}) []interface{} {
	items := make([]interface{}, 0)
	for ele := range c {
		items = append(items, ele)
	}
	return items
}

func TestInOrderItems(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	val.Sample("abc")
	val.Sample("abc")
	val.Sample("test")
	val.Sample("abc")
	val.Sample("qq")

	items := val.ItemsTop(2)

	assert.Equal(t, 2, len(items), "Expected top 2")
	assert.Equal(t, "abc", items[0].Name)
	assert.Equal(t, "test", items[1].Name)
}

func TestIteration(t *testing.T) {
	val := NewCounter()
	val.Sample("a")
	val.Sample("b")
	val.Sample("c")

	total := 0
	for range val.Iter() {
		total++
	}

	assert.Equal(t, 3, total)
}
