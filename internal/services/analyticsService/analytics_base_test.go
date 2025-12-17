package analyticsService

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_ReturnsService(t *testing.T) {
	s := New(nil, nil)
	require.NotNil(t, s)
}

func TestGetTop_InvalidMetric_Error(t *testing.T) {
	s := New(nil, nil)
	out, err := s.GetTop(context.Background(), "invalid", 1)
	require.ErrorIs(t, err, ErrInvalidMetric)
	require.Nil(t, out)
}

func TestGetAuthorStats_ErrorFromStorage(t *testing.T) {
	store := NewMockStorage(t)
	cache := NewMockCache(t)
	svc := New(store, cache)

	store.EXPECT().
		GetAuthorStats(context.Background(), "42").
		Return(nil, errors.New("pg fail"))

	out, err := svc.GetAuthorStats(context.Background(), "42")
	require.Error(t, err)
	require.Nil(t, out)
}
