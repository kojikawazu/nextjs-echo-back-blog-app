package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogById(t *testing.T) {
	// モックデータ
	mockBlogData := &models.BlogData{
		ID:         testBlogID,
		BlogUserId: testUserID,
		Title:      "title1",
		GithubUrl:  "https://github.com/user/repo1",
		Category:   "Category1",
		Tags:       "Tag1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name              string
		id                string
		setupMock         func(m *repositories_blogs.MockBlogRepository)
		wantBlog          *models.BlogData
		wantErr           error // 期待する sentinel（nil は正常系）
		wantInternal      bool  // 想定外エラー（sentinel に変換されず原因を保持する）
		expectRepoNotCall bool  // バリデーションでリポジトリが呼ばれないケース
	}{
		{
			name: "正常系_ブログが存在する場合はブログを返す",
			id:   testBlogID,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogById", testBlogID).Return(mockBlogData, nil)
			},
			wantBlog: mockBlogData,
		},
		{
			name:              "準正常系_idが空文字の場合ErrInvalidIDを返す",
			id:                "",
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidID,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_idがUUID形式でない場合ErrInvalidIDを返す",
			id:                "1",
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidID,
			expectRepoNotCall: true,
		},
		{
			name: "準正常系_ブログが存在しない場合ErrBlogNotFoundを返す",
			id:   testBlogID,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogById", testBlogID).Return(nil, repositories_blogs.ErrBlogNotFound)
			},
			wantErr: services_blogs.ErrBlogNotFound,
		},
		{
			name: "異常系_Repositoryが想定外エラーを返す場合はErrBlogNotFoundにしない",
			id:   testBlogID,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogById", testBlogID).Return(nil, errRepositoryFailure)
			},
			wantInternal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)
			tt.setupMock(mockRepo)

			blog, err := blogService.FetchBlogById(tt.id)

			switch {
			case tt.wantInternal:
				assertInternalError(t, err, errRepositoryFailure)
				assert.NotErrorIs(t, err, services_blogs.ErrBlogNotFound)
				assert.Nil(t, blog)
			case tt.wantErr != nil:
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, blog)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBlog, blog)
			}
			if tt.expectRepoNotCall {
				mockRepo.AssertNotCalled(t, "FetchBlogById", tt.id)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
