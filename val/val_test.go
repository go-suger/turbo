package val

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	str := "dada"
	ptr := Ptr[string](str)
	assert.NotNil(t, ptr)
	assert.Equal(t, str, *ptr)
}
