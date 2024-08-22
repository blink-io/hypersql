package hypersql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualFold(t *testing.T) {
	var str1 = "Hello"
	var str2 = "HELLo"
	var flag = strings.EqualFold(str1, str2)

	assert.True(t, flag)
}
