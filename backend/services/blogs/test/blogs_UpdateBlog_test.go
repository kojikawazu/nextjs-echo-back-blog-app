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

func TestService_UpdateBlog(t *testing.T) {
	// 正常な入力データ
	validInput := struct {
		id, title, githubURL, category, description, tags string
	}{
		id: testBlogID, title: "Test Blog", githubURL: "https://github.com/user/repo",
		category: "Tech", description: "This is a test blog.", tags: "go, testing",
	}

	// 期待されるブログデータ
	expectedBlog := &models.BlogData{
		ID:          validInput.id,
		BlogUserId:  testUserID,
		Title:       validInput.title,
		GithubUrl:   validInput.githubURL,
		Category:    validInput.category,
		Description: validInput.description,
		Tags:        validInput.tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// withInput は正常な入力の一部を差し替えた入力を返す
	type input struct{ id, title, githubURL, category, description, tags string }
	base := input(validInput)
	withInput := func(f func(in *input)) input {
		in := base
		f(&in)
		return in
	}

	// expectUpdate は正常な入力で UpdateBlog が呼ばれるモックを設定する
	expectUpdate := func(blog *models.BlogData, err error) func(m *repositories_blogs.MockBlogRepository) {
		return func(m *repositories_blogs.MockBlogRepository) {
			m.On("UpdateBlog", base.id, base.title, base.githubURL, base.category, base.description, base.tags).Return(blog, err)
		}
	}

	tests := []struct {
		name              string
		in                input
		setupMock         func(m *repositories_blogs.MockBlogRepository)
		wantErr           error // 期待する sentinel（nil は正常系）
		wantInternal      bool  // 想定外エラー（sentinel に変換されず原因を保持する）
		expectRepoNotCall bool  // バリデーションでリポジトリが呼ばれないケース
	}{
		{
			name:      "正常系_全フィールド正常値でブログ更新成功",
			in:        base,
			setupMock: expectUpdate(expectedBlog, nil),
		},
		{
			name:              "準正常系_idが空文字の場合ErrInvalidIDを返す",
			in:                withInput(func(in *input) { in.id = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidID,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_idがUUID形式でない場合ErrInvalidIDを返す",
			in:                withInput(func(in *input) { in.id = "123" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidID,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_titleが空文字の場合ErrInvalidTitleを返す",
			in:                withInput(func(in *input) { in.title = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidTitle,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_githubUrlが空文字の場合ErrInvalidGithubURLを返す",
			in:                withInput(func(in *input) { in.githubURL = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidGithubURL,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_categoryが空文字の場合ErrInvalidCategoryを返す",
			in:                withInput(func(in *input) { in.category = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidCategory,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_descriptionが空文字の場合ErrInvalidDescriptionを返す",
			in:                withInput(func(in *input) { in.description = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidDescription,
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_tagsが空文字の場合ErrInvalidTagsを返す",
			in:                withInput(func(in *input) { in.tags = "" }),
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           services_blogs.ErrInvalidTags,
			expectRepoNotCall: true,
		},
		{
			name:      "準正常系_ブログが存在しない場合ErrBlogNotFoundを返す",
			in:        base,
			setupMock: expectUpdate(nil, repositories_blogs.ErrBlogNotFound),
			wantErr:   services_blogs.ErrBlogNotFound,
		},
		{
			name:         "異常系_Repositoryが想定外エラーを返す場合はErrBlogNotFoundにしない",
			in:           base,
			setupMock:    expectUpdate(nil, errRepositoryFailure),
			wantInternal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)
			tt.setupMock(mockRepo)

			blog, err := blogService.UpdateBlog(tt.in.id, tt.in.title, tt.in.githubURL, tt.in.category, tt.in.description, tt.in.tags)

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
				assert.Equal(t, expectedBlog, blog)
			}
			if tt.expectRepoNotCall {
				mockRepo.AssertNotCalled(t, "UpdateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
