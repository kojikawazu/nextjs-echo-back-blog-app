package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogs(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		setupMock  func(m *repositories_blogs.MockBlogRepository)
		wantLen    int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常系_ブログ一覧を複数件取得成功",
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogs").Return([]models.BlogData{
					{ID: "1", BlogUserId: "1", Title: "title1", GithubUrl: "https://github.com/user/repo1", Category: "Category1", Tags: "Tag1", CreatedAt: now, UpdatedAt: now},
					{ID: "2", BlogUserId: "2", Title: "title2", GithubUrl: "https://github.com/user/repo2", Category: "Category2", Tags: "Tag2", CreatedAt: now, UpdatedAt: now},
				}, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "正常系_ブログが0件の場合空スライスを返す",
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogs").Return([]models.BlogData{}, nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "異常系_Repositoryがエラーを返す場合エラーを返す",
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogs").Return(nil, errors.New("db connection error"))
			},
			wantErr:    true,
			wantErrMsg: "db connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)

			tt.setupMock(mockRepo)

			blogs, err := blogService.FetchBlogs()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Equal(t, tt.wantErrMsg, err.Error())
				}
				assert.Nil(t, blogs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, blogs)
				assert.Len(t, blogs, tt.wantLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
