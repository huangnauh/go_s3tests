package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	assert := assert.New(t)

	assert.Equal(LoadConfig(), nil)
}
