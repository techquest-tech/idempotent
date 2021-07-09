package idempotent_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/techquest-tech/idempotent"
)

func TestKey(t *testing.T) {
	obj := map[string]interface{}{
		"ID":  "Testing",
		"Now": time.Now(),
	}
	idens, err := idempotent.TemplateAsKey(`{{ printf "%s-%s" .ID (.Now | FormateDate) }}`)
	assert.Nil(t, err)
	// t.Log("key:", keyvalue)

	idens.Target = obj
	key, err := idens.IdempotentKey()
	assert.Nil(t, err)
	log.Print("key=", key)
}
