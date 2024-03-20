package testutils_test

import (
	"github.com/stretchr/testify/require"
	"go-labs/internal/testutils"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	addr, cleanupFunc := testutils.NewRedisClient()
	defer func() {
		require.NoError(t, cleanupFunc())
	}()
	require.NotNil(t, addr)
}
