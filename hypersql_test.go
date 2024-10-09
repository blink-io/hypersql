package hypersql

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEqualFold(t *testing.T) {
	var str1 = "Hello"
	var str2 = "HELLo"
	var flag = strings.EqualFold(str1, str2)

	assert.True(t, flag)
}

func TestLocation_1(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	require.NotNil(t, loc)
}
