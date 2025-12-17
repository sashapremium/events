package analyticsService

import (
	"context"
	"errors"
	"testing"

	"github.com/sashapremium/events/internal/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IngestSuite struct {
	suite.Suite

	ctx   context.Context
	svc   *Service
	store *MockStorage
	cache *MockCache
}

func TestIngestSuite(t *testing.T) {
	suite.Run(t, new(IngestSuite))
}

func (s *IngestSuite) SetupTest() {
	s.ctx = context.Background()
	s.store = NewMockStorage(s.T())
	s.cache = NewMockCache(s.T())
	s.svc = New(s.store, s.cache)
}

func (s *IngestSuite) TestProcessEvent_NilEvent_NoError() {
	err := s.svc.ProcessEvent(s.ctx, nil)
	require.NoError(s.T(), err)
}

func (s *IngestSuite) TestProcessEvent_InvalidContentID_ReturnsError() {
	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "not-a-number",
		Type:      models.EventView,
	})
	require.Error(s.T(), err)
}

func (s *IngestSuite) TestProcessEvent_UnknownType_NoError_NoCalls() {
	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "10",
		Type:      "UNKNOWN",
	})
	require.NoError(s.T(), err)
}

func (s *IngestSuite) TestProcessEvent_View_WithUniqueUser_Added() {
	postID := uint64(10)

	s.cache.EXPECT().
		AddUniqueUser(s.ctx, postID, "u1").
		Return(true, nil)

	s.cache.EXPECT().
		IncrTotals(s.ctx, postID, TotalsDelta{Views: 1, UniqueUsers: 1}).
		Return(nil)

	s.cache.EXPECT().
		MarkDirty(s.ctx, postID).
		Return(nil)

	s.cache.EXPECT().
		IncrTop(s.ctx, "views", postID, int64(1)).
		Return(nil)

	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "10",
		Type:      models.EventView,
		UserHash:  "u1",
	})

	require.NoError(s.T(), err)
}

func (s *IngestSuite) TestProcessEvent_View_AddUniqueUserError_Stops() {
	postID := uint64(10)
	wantErr := errors.New("redis down")

	s.cache.EXPECT().
		AddUniqueUser(s.ctx, postID, "u1").
		Return(false, wantErr)

	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "10",
		Type:      models.EventView,
		UserHash:  "u1",
	})

	require.ErrorIs(s.T(), err, wantErr)
}

func (s *IngestSuite) TestProcessEvent_Like_IncrTotalsError_Stops() {
	postID := uint64(11)
	wantErr := errors.New("incr totals failed")

	s.cache.EXPECT().
		IncrTotals(s.ctx, postID, TotalsDelta{Likes: 1}).
		Return(wantErr)

	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "11",
		Type:      models.EventLike,
	})

	require.ErrorIs(s.T(), err, wantErr)
}

func (s *IngestSuite) TestProcessEvent_Like_MarkDirtyError_Stops() {
	postID := uint64(12)
	wantErr := errors.New("mark dirty failed")

	s.cache.EXPECT().
		IncrTotals(s.ctx, postID, TotalsDelta{Likes: 1}).
		Return(nil)

	s.cache.EXPECT().
		MarkDirty(s.ctx, postID).
		Return(wantErr)

	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "12",
		Type:      models.EventLike,
	})

	require.ErrorIs(s.T(), err, wantErr)
}

func (s *IngestSuite) TestProcessEvent_Like_IncrTopError_Stops() {
	postID := uint64(13)
	wantErr := errors.New("top failed")

	s.cache.EXPECT().
		IncrTotals(s.ctx, postID, TotalsDelta{Likes: 1}).
		Return(nil)

	s.cache.EXPECT().
		MarkDirty(s.ctx, postID).
		Return(nil)

	s.cache.EXPECT().
		IncrTop(s.ctx, "likes", postID, int64(1)).
		Return(wantErr)

	err := s.svc.ProcessEvent(s.ctx, &models.ContentEvent{
		ContentID: "13",
		Type:      models.EventLike,
	})

	require.ErrorIs(s.T(), err, wantErr)
}
