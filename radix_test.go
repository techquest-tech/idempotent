package idempotent_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/techquest-tech/idempotent"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateData() {
	// r := iradix.New()
	rt := idempotent.NewRadixTree()

	test := idempotent.PersistedIdempotent{
		File:     "b.tmp",
		Service:  rt,
		Interval: "5s",
	}
	test.Init()

	for i := 0; i < 100000; i++ {
		txt := RandStringRunes(64)
		test.Save(txt)

	}

	time.Sleep(29 * time.Second)

}

func TestKeys(t *testing.T) {
	rt := idempotent.NewRadixTree()

	test := idempotent.PersistedIdempotent{
		File:     "b.tmp",
		Service:  rt,
		Interval: "5s",
	}
	test.Init()

	result, err := test.Duplicated("kzutEpNCwjvqzSoDWNLisIdUSShCxJjeXExRHzoZcZfRIznoxRvytamxmxhTttQL")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = test.Duplicated("hello Radix")
	assert.Nil(t, err)
	assert.False(t, result)

	// keys, err := test.AllKeys()
	// assert.Nil(t, err)
	// for _, item := range keys {
	// 	fmt.Println(item)
	// }
}
