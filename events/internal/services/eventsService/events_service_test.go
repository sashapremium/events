package eventsService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	eventmodel "github.com/sashapremium/events/events/internal/models"
	"github.com/sashapremium/events/events/internal/services/eventsService/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx      context.Context
	storage  *mocks.StorageMock
	producer *mocks.ProducerMock
	service  *Service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = &mocks.StorageMock{}
	s.producer = &mocks.ProducerMock{}
	s.service = NewService(s.storage, s.producer)
	s.storage.On("GetPostAuthorID", mock.Anything, mock.AnythingOfType("uint64")).
		Return(uint64(1), nil).
		Maybe()

	s.storage.On("GetPost", mock.Anything, mock.AnythingOfType("uint64")).
		Return((*eventmodel.PostInfo)(nil), fmt.Errorf("GetPost: not implemented")).
		Maybe()
}
func (s *ServiceSuite) TestGetPost_NotImplemented() {
	_, err := s.service.GetPost(s.ctx, 1)
	s.Require().Error(err)
}
func (s *ServiceSuite) TestViewPost_Success() {
	postID := uint64(40410)
	userHash := "user-123"
	wantContentID := strconv.FormatUint(postID, 10)
	s.storage.
		On("InsertEvents", mock.Anything, mock.MatchedBy(func(events []*eventmodel.ContentEvent) bool {
			s.Require().Len(events, 1)
			ev := events[0]

			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventView, ev.Type)
			s.False(ev.At.IsZero())

			return true
		})).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.MatchedBy(func(ev *eventmodel.ContentEvent) bool {
			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventView, ev.Type)
			return true
		})).
		Return(nil).
		Once()

	err := s.service.ViewPost(s.ctx, postID, userHash)
	s.Require().NoError(err)
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLikePost_Success() {
	postID := uint64(40410)
	userHash := "user-123"
	wantContentID := strconv.FormatUint(postID, 10)
	s.storage.
		On("InsertEvents", mock.Anything, mock.MatchedBy(func(events []*eventmodel.ContentEvent) bool {
			s.Require().Len(events, 1)
			ev := events[0]

			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventLike, ev.Type)
			s.False(ev.At.IsZero())

			return true
		})).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.MatchedBy(func(ev *eventmodel.ContentEvent) bool {
			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventLike, ev.Type)
			return true
		})).
		Return(nil).
		Once()

	err := s.service.LikePost(s.ctx, postID, userHash)
	s.Require().NoError(err)

	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLikePost_ValidationError() {
	err := s.service.LikePost(s.ctx, 1, "")
	s.Require().Error(err)

	s.storage.AssertNotCalled(s.T(), "InsertEvents", mock.Anything, mock.Anything)
	s.producer.AssertNotCalled(s.T(), "PublishEvent", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestLikePost_StorageError() {
	postID := uint64(1)
	userHash := "u"
	wantErr := errors.New("storage error")

	s.storage.
		On("InsertEvents", mock.Anything, mock.Anything).
		Return(wantErr).
		Once()

	err := s.service.LikePost(s.ctx, postID, userHash)

	s.Require().Error(err)
	s.ErrorIs(err, wantErr)

	s.producer.AssertNotCalled(s.T(), "PublishEvent", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestLikePost_ProducerError() {
	postID := uint64(10)
	userHash := "u"
	wantErr := errors.New("producer error")

	s.storage.
		On("InsertEvents", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.Anything).
		Return(wantErr).
		Once()

	err := s.service.LikePost(s.ctx, postID, userHash)

	s.Require().Error(err)
	s.ErrorIs(err, wantErr)
}

func (s *ServiceSuite) TestRepostPost_Success() {
	postID := uint64(40410)
	userHash := "user-123"
	wantContentID := strconv.FormatUint(postID, 10)

	s.storage.
		On("InsertEvents", mock.Anything, mock.MatchedBy(func(events []*eventmodel.ContentEvent) bool {
			s.Require().Len(events, 1)
			ev := events[0]

			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventRepost, ev.Type)

			return true
		})).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.MatchedBy(func(ev *eventmodel.ContentEvent) bool {
			s.Equal(eventmodel.EventRepost, ev.Type)
			return true
		})).
		Return(nil).
		Once()

	err := s.service.RepostPost(s.ctx, postID, userHash)
	s.Require().NoError(err)

	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRepostPost_ValidationError() {
	err := s.service.RepostPost(s.ctx, 1, "")
	s.Require().Error(err)

	s.storage.AssertNotCalled(s.T(), "InsertEvents", mock.Anything, mock.Anything)
	s.producer.AssertNotCalled(s.T(), "PublishEvent", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestAddComment_Success() {
	postID := uint64(40410)
	userHash := "user-123"
	text := "Интересный пост"
	wantContentID := strconv.FormatUint(postID, 10)

	before := time.Now()

	s.storage.
		On("InsertEvents", mock.Anything, mock.MatchedBy(func(events []*eventmodel.ContentEvent) bool {
			s.Require().Len(events, 1)
			ev := events[0]

			s.Equal(wantContentID, ev.ContentID)
			s.Equal(userHash, ev.UserHash)
			s.Equal(eventmodel.EventComment, ev.Type)
			s.Equal(text, ev.Comment)
			s.True(ev.At.After(before) || ev.At.Equal(before))

			return true
		})).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.MatchedBy(func(ev *eventmodel.ContentEvent) bool {
			s.Equal(eventmodel.EventComment, ev.Type)
			s.Equal(text, ev.Comment)
			return true
		})).
		Return(nil).
		Once()

	comment, err := s.service.AddComment(s.ctx, postID, userHash, text)

	s.Require().NoError(err)
	s.Require().NotNil(comment)

	s.Equal(postID, comment.PostId)
	s.Equal(userHash, comment.UserHash)
	s.Equal(text, comment.Text)
	s.NotEmpty(comment.CreatedAt)

	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAddComment_EmptyTextValidationError() {
	comment, err := s.service.AddComment(s.ctx, 1, "u", "   ")
	s.Require().Error(err)
	s.Nil(comment)

	s.storage.AssertNotCalled(s.T(), "InsertEvents", mock.Anything, mock.Anything)
	s.producer.AssertNotCalled(s.T(), "PublishEvent", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestAddComment_TooLongValidationError() {
	text := strings.Repeat("a", 1001) //  лимит 1000
	comment, err := s.service.AddComment(s.ctx, 1, "u", text)
	s.Require().Error(err)
	s.Nil(comment)

	s.storage.AssertNotCalled(s.T(), "InsertEvents", mock.Anything, mock.Anything)
	s.producer.AssertNotCalled(s.T(), "PublishEvent", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestAddComment_ProducerError() {
	postID := uint64(1)
	userHash := "u"
	text := "ok"
	wantErr := errors.New("producer error")

	s.storage.
		On("InsertEvents", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	s.producer.
		On("PublishEvent", mock.Anything, mock.Anything).
		Return(wantErr).
		Once()

	comment, err := s.service.AddComment(s.ctx, postID, userHash, text)
	s.Require().Error(err)
	s.ErrorIs(err, wantErr)
	s.Nil(comment)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
func (s *ServiceSuite) TestPersistAndPublish_StorageNil() {
	s.service = NewService(nil, s.producer)

	err := s.service.LikePost(s.ctx, 1, "u")
	s.Require().Error(err)
}

func (s *ServiceSuite) TestPersistAndPublish_ProducerNil() {
	s.service = NewService(s.storage, nil)

	err := s.service.LikePost(s.ctx, 1, "u")
	s.Require().Error(err)
}
