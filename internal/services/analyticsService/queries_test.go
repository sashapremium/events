package analyticsService

import (
	"context"
	"errors"
	"testing"
	"time"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueriesSuite struct {
	suite.Suite

	ctx   context.Context
	svc   *Service
	store *MockStorage
	cache *MockCache
}

func TestQueriesSuite(t *testing.T) {
	suite.Run(t, new(QueriesSuite))
}

func (s *QueriesSuite) SetupTest() {
	s.ctx = context.Background()

	// ВАЖНО:
	// Моки сгенерированы в package analyticsService (outpkg: analyticsService),
	// поэтому никаких imports ".../mocks" и никакого префикса mocks. тут нет.
	s.store = NewMockStorage(s.T())
	s.cache = NewMockCache(s.T())

	s.svc = New(s.store, s.cache)
}

func (s *QueriesSuite) TestGetPostStats_StorageError() {
	postID := uint64(10)

	s.store.EXPECT().
		GetPostTotals(s.ctx, postID).
		Return(PostTotals{}, errors.New("db down"))

	out, err := s.svc.GetPostStats(s.ctx, postID, false)
	require.Error(s.T(), err)
	require.Nil(s.T(), out)
}

func (s *QueriesSuite) TestGetPostStats_FreshFalse_LastSyncedOk() {
	postID := uint64(10)

	s.store.EXPECT().
		GetPostTotals(s.ctx, postID).
		Return(PostTotals{
			Views:       100,
			Likes:       5,
			Comments:    2,
			Reposts:     1,
			UniqueUsers: 90,
		}, nil)

	s.cache.EXPECT().
		GetLastSyncedAt(s.ctx).
		Return("2025-01-01T00:00:00Z", true, nil)

	out, err := s.svc.GetPostStats(s.ctx, postID, false)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), out)

	require.Equal(s.T(), postID, out.PostId)
	require.NotNil(s.T(), out.Totals)

	require.EqualValues(s.T(), 100, out.Totals.TotalViews)
	require.EqualValues(s.T(), 5, out.Totals.TotalLikes)
	require.EqualValues(s.T(), 2, out.Totals.TotalComments)
	require.EqualValues(s.T(), 1, out.Totals.TotalReposts)
	require.EqualValues(s.T(), 90, out.Totals.UniqueUsers)

	require.Equal(s.T(), "2025-01-01T00:00:00Z", out.LastSyncedAt)
	require.Nil(s.T(), out.FreshTail)
}

func (s *QueriesSuite) TestGetPostStats_LastSyncedError_FallbackZeroTime() {
	postID := uint64(11)

	s.store.EXPECT().
		GetPostTotals(s.ctx, postID).
		Return(PostTotals{}, nil)

	s.cache.EXPECT().
		GetLastSyncedAt(s.ctx).
		Return("", false, errors.New("redis err"))

	out, err := s.svc.GetPostStats(s.ctx, postID, false)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), out)

	require.Equal(s.T(), time.Time{}.UTC().Format(time.RFC3339), out.LastSyncedAt)
}

func (s *QueriesSuite) TestGetPostStats_FreshTrue_DeltaHit() {
	postID := uint64(12)

	s.store.EXPECT().
		GetPostTotals(s.ctx, postID).
		Return(PostTotals{}, nil)

	s.cache.EXPECT().
		GetLastSyncedAt(s.ctx).
		Return("2025-01-01T00:00:00Z", true, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(TotalsDelta{
			Views:       7,
			Likes:       1,
			Comments:    2,
			Reposts:     0,
			UniqueUsers: 3,
		}, true, nil)

	out, err := s.svc.GetPostStats(s.ctx, postID, true)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), out)
	require.NotNil(s.T(), out.FreshTail)

	require.EqualValues(s.T(), 7, out.FreshTail.Views)
	require.EqualValues(s.T(), 1, out.FreshTail.Likes)
	require.EqualValues(s.T(), 2, out.FreshTail.Comments)
	require.EqualValues(s.T(), 0, out.FreshTail.Reposts)
	require.EqualValues(s.T(), 3, out.FreshTail.UniqueUsers)
}

func (s *QueriesSuite) TestGetTop_Success() {
	metric := "views"
	limit := uint32(3)

	s.cache.EXPECT().
		GetTop(s.ctx, metric, limit).
		Return([]TopItem{
			{PostID: 1, Value: 100},
			{PostID: 2, Value: 90},
		}, nil)

	out, err := s.svc.GetTop(s.ctx, metric, limit)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), out)

	require.Equal(s.T(), metric, out.Metric)
	require.Len(s.T(), out.Items, 2)

	require.EqualValues(s.T(), 1, out.Items[0].PostId)
	require.EqualValues(s.T(), 100, out.Items[0].Value)
}

func (s *QueriesSuite) TestGetAuthorStats_Passthrough() {
	authorID := "author-1"
	want := &analyticsmodels.AuthorStatsModel{
		AuthorId:      authorID,
		Posts:         7,
		TotalViews:    100,
		TotalLikes:    10,
		TotalComments: 3,
		TotalReposts:  1,
	}

	s.store.EXPECT().
		GetAuthorStats(s.ctx, authorID).
		Return(want, nil)

	out, err := s.svc.GetAuthorStats(s.ctx, authorID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, out)
}
