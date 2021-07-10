package idempotent_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techquest-tech/idempotent"
)

func TestToByte(t *testing.T) {

	idempotent.ToBytes("abcd")

	key := uint32(rand.Intn(321189))
	result, err := idempotent.ToBytes(key)
	assert.Nil(t, err)
	fmt.Println(string(result))
	fmt.Printf("key: %d\n", key)
	for _, item := range result {
		fmt.Printf("%d\n", item)
	}
}

func BenchmarkToByte(b *testing.B) {

	for i := 0; i < b.N; i++ {
		idempotent.ToBytes(rand.Intn(10000))
	}
}
