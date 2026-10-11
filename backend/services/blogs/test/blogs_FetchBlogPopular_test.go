package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_FetchBlogPopular(t *testing.T) {
	mockBlogData := []models.BlogData{
		{
			ID:         "1",
			BlogUserId: "1",
			Title:      "title1",
			GithubUrl:  "https://github.com/user/repo1",
			Category:   "Category1",
			Tags:       "Tag1",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "2",
			BlogUserId: "2",
			Title:      "title2",
			GithubUrl:  "https://github.com/user/repo2",
			Category:   "Category2",
			Tags:       "Tag2",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	tests := []struct {
		name              string
		count             int
		setupMock         func(m *repositories_blogs.MockBlogRepository)
		wantLen           int
		wantErr           error // 期待する sentinel（nil は正常系）
		wantInternal      bool  // 想定外エラー（sentinel に変換されず原因を保持する）
		expectRepoNotCall bool  // バリデーションでリポジトリが呼ばれないケース
	}{
		{
			name:  "正常系_指定件数の人気ブログを返す",
			count: 2,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogPopular", 2).Return(mockBlogData, nil)
			},
			wantLen: 2,
		},
		{
			name:              "準正常系_countが0の場合ErrInvalidCountを返す",
			count:             0,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidCount,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_countが負数の場合ErrInvalidCountを返す",
			count:             -1,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidCount,
			expectRepoNotCall: true,
		},
		{
			name:  "異常系_Repositoryがエラーを返す場合は原因を保持した想定外エラーを返す",
			count: 2,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("FetchBlogPopular", 2).Return(nil, errRepositoryFailure)
			},
			wantInternal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)
			tt.setupMock(mockRepo)

			blogs, err := blogService.FetchBlogPopular(tt.count)

			switch {
			case tt.wantInternal:
				assertInternalError(t, err, errRepositoryFailure)
				assert.Nil(t, blogs)
			case tt.wantErr != nil:
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, blogs)
			default:
				assert.NoError(t, err)
				assert.Len(t, blogs, tt.wantLen)
			}
			if tt.expectRepoNotCall {
				mockRepo.AssertNotCalled(t, "FetchBlogPopular", mock.Anything)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
