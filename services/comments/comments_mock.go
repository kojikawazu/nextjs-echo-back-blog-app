package services_comments

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) FetchCommentsByBlogId(blogId string) ([]models.CommentData, error) {
	args := m.Called(blogId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CommentData), args.Error(1)
}

func (m *MockCommentService) CreateComment(blogId, guestUser, comment string) (*models.CommentData, error) {
	args := m.Called(blogId, guestUser, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CommentData), args.Error(1)
}
