package bolt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextSequence(t *testing.T) {
	boltClient := NewBoltClient("short-url.db", 0600)
	err := boltClient.InitialBucket()
	assert.Nil(t, err)
	for i := 1; i <= 100; i++ {
		seq, err := boltClient.NextSequence()
		assert.Nil(t, err)
		assert.EqualValues(t, int(seq), i)
	}
}

func TestForEach(t *testing.T) {
	boltClient := NewBoltClient("short-url.db", 0600)
	err := boltClient.InitialBucket()
	assert.Nil(t, err)
	result, err := boltClient.ForEach()
	for key, value := range result {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func TestSetAndGetValue(t *testing.T) {
	boltClient := NewBoltClient("short-url.db", 0600)
	err := boltClient.InitialBucket()
	assert.Nil(t, err)
	boltClient.Set("sag/2$#", "1")
	boltClient.Set("!@##$", "2")
	boltClient.Set("+_()*(", "3")
	v1, err := boltClient.Get("sag/2$#")
	assert.Nil(t, err)
	assert.EqualValues(t, v1, []byte("1"))
	v2, err := boltClient.Get("!@##$")
	assert.Nil(t, err)
	assert.EqualValues(t, string(v2), "2")
	v3, err := boltClient.Get("+_()*(")
	assert.Nil(t, err)
	assert.EqualValues(t, string(v3), "3")
}
