package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartup(t *testing.T) {
	err := startup(context.Background(), "production")
	assert.NoError(t, err)
}
