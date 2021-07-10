package idempotent_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/techquest-tech/idempotent"
)

type IdTest struct {
	ID string
}

func (id IdTest) IdempotentKey() (interface{}, error) {
	return id.ID, nil
}

func TestDuplicated(t *testing.T) {

	obj := map[string]interface{}{
		"ID":  "Testing",
		"Now": time.Now(),
	}

	ram := idempotent.NewInMemoryMap()

	factory, err := idempotent.NewIdempotentWithTemplate(`{{ printf "%s-%s" .ID (.Now | FormateDate) }}`, ram)
	assert.Nil(t, err)

	result, err := factory.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, !result)

	result, err = factory.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, result)

	obj["ID"] = "updated ID"
	result, err = factory.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, !result)
}

func TestDuplicatedFile(t *testing.T) {
	os.Remove("test.txt")
	obj := IdTest{
		ID: "TestDuplicatedFile",
	}
	idd := idempotent.PersistedIdempotent{
		File:    "test.txt",
		Service: idempotent.NewInMemoryMap(),
	}
	err := idd.Init()
	assert.Nil(t, err)

	factory := &idempotent.Idempotent{
		Service: &idd,
	}
	result, err := factory.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, !result)

	result, err = factory.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, result)

	os.Remove("test.txt")
	all, _ := factory.Service.AllKeys()
	fmt.Printf("%+v", all)
}

type TestMap map[string]interface{}

func TestMapDuplicated(t *testing.T) {

	maptesting, err := idempotent.NewIdempotentWithTemplate(`{{ .ID }}`, idempotent.NewInMemoryMap())

	assert.Nil(t, err)

	obj := map[string]interface{}{
		"ID":   9999,
		"body": "hello world",
	}

	result, err := maptesting.Duplicated(obj)
	assert.Nil(t, err)
	assert.False(t, result)

	obj["other attribue"] = time.Now()

	result, err = maptesting.Duplicated(obj)
	assert.Nil(t, err)
	assert.True(t, result)

	obj["ID"] = 10000

	result, err = maptesting.Duplicated(obj)
	assert.Nil(t, err)
	assert.False(t, result)

	obj2 := TestMap{
		"ID":   9999,
		"Body": "it's another type",
	}

	result, err = maptesting.Duplicated(obj2)
	assert.Nil(t, err)
	assert.True(t, result)

	all, _ := maptesting.Service.AllKeys()
	fmt.Printf("%+v\n", all)

}

func BenchmarkIdempotent(b *testing.B) {
	factory := &idempotent.Idempotent{
		Service: idempotent.NewInMemoryMap(),
	}

	for i := 0; i < b.N; i++ {
		rnumber := RandStringRunes(64)
		dup, err := factory.Duplicated(rnumber)
		assert.Nil(b, err)
		assert.False(b, dup)
	}

}

func BenchmarkRadix(b *testing.B) {
	factory := &idempotent.Idempotent{
		Service: idempotent.NewRadixTree(),
	}

	for i := 0; i < b.N; i++ {
		rnumber := RandStringRunes(64)
		dup, err := factory.Duplicated(rnumber)
		assert.Nil(b, err)
		assert.False(b, dup)
	}

}
