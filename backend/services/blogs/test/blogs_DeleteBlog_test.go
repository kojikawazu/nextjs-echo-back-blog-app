package services_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_DeleteBlog(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		setupMock         func(m *repositories_blogs.MockBlogRepository)
		wantErr           error // 期待する sentinel（nil は正常系）
		wantInternal      bool  // 想定外エラー（sentinel に変換されず原因を保持する）
		expectRepoNotCall bool  // バリデーションでリポジトリが呼ばれないケース
	}{
		{
			name: "正常系_UUID形式のidで削除成功",
			id:   testBlogID,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("DeleteBlog", testBlogID).Return(nil)
			},
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
			id:                "123",
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidID,
			expectRepoNotCall: true,
		},
		{
			name: "異常系_Repositoryがエラーを返す場合は原因を保持した想定外エラーを返す",
			id:   testBlogID,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("DeleteBlog", testBlogID).Return(errRepositoryFailure)
			},
			wantInternal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)
			tt.setupMock(mockRepo)

			err := blogService.DeleteBlog(tt.id)

			switch {
			case tt.wantInternal:
				assertInternalError(t, err, errRepositoryFailure)
			case tt.wantErr != nil:
				assert.ErrorIs(t, err, tt.wantErr)
			default:
				assert.NoError(t, err)
			}
			if tt.expectRepoNotCall {
				mockRepo.AssertNotCalled(t, "DeleteBlog", tt.id)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
