package repositories_blog_comments

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error) {
	args := m.Called(blogId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogCommentsData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentRepository) CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error) {
	args := m.Called(blogId, guestUser, comment)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogCommentsData), args.Error(1)
	}
	return nil, args.Error(1)
}
