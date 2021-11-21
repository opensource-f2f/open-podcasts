package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRSSSources(t *testing.T) {
	items, err := GetRSSSources()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(items))
}
