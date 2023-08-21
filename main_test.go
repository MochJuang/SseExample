package main_test

import (
	code "go-sse"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	t.Run("negative test", func(t *testing.T) {
		result := code.Sum(2, 2)
		assert.NotEqual(t, result, 5)
	})

	t.Run("positive test", func(t *testing.T) {
		result := code.Sum(2, 2)
		assert.Equal(t, result, 4)
	})
}
