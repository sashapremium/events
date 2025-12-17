package analyticsService

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMetric_AllBranches(t *testing.T) {
	store := NewMockStorage(t)
	cache := NewMockCache(t)
	svc := New(store, cache)

	valid := []string{"views", "likes", "comments", "reposts"}
	for _, m := range valid {
		err := svc.validateMetric(m)
		require.NoError(t, err)
	}

	err := svc.validateMetric("unknown")
	require.ErrorIs(t, err, ErrInvalidMetric)
}

func TestGetTop_ValidMetric_CallsCache(t *testing.T) {
	store := NewMockStorage(t)
	cache := NewMockCache(t)
	svc := New(store, cache)

	cache.EXPECT().
		GetTop(context.Background(), "views", uint32(5)).
		Return([]TopItem{{PostID: 1, Value: 100}}, nil)

	out, err := svc.GetTop(context.Background(), "views", 5)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)
}
