package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateBlog(t *testing.T) {
	validInput := struct {
		userId, title, githubURL, category, description, tags string
	}{
		userId: "user1", title: "Test Blog", githubURL: "https://github.com/user/repo",
		category: "Tech", description: "This is a test blog.", tags: "go, testing",
	}

	expectedBlog := models.BlogData{
		ID:          "123",
		BlogUserId:  validInput.userId,
		Title:       validInput.title,
		GithubUrl:   validInput.githubURL,
		Category:    validInput.category,
		Description: validInput.description,
		Tags:        validInput.tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name              string
		userId            string
		title             string
		githubURL         string
		category          string
		description       string
		tags              string
		setupMock         func(m *repositories_blogs.MockBlogRepository)
		wantErr           bool
		wantErrMsg        string
		expectRepoNotCall bool // バリデーションでリポジトリが呼ばれないケース
	}{
		{
			name:        "正常系_全フィールド正常値でブログ作成成功",
			userId:      validInput.userId,
			title:       validInput.title,
			githubURL:   validInput.githubURL,
			category:    validInput.category,
			description: validInput.description,
			tags:        validInput.tags,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("CreateBlog", validInput.userId, validInput.title, validInput.githubURL, validInput.category, validInput.description, validInput.tags).Return(&expectedBlog, nil)
			},
			wantErr: false,
		},
		{
			name:              "準正常系_userIdが空文字の場合エラーを返す",
			userId:            "",
			title:             validInput.title,
			githubURL:         validInput.githubURL,
			category:          validInput.category,
			description:       validInput.description,
			tags:              validInput.tags,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid userId",
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_titleが空文字の場合エラーを返す",
			userId:            validInput.userId,
			title:             "",
			githubURL:         validInput.githubURL,
			category:          validInput.category,
			description:       validInput.description,
			tags:              validInput.tags,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid title",
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_githubUrlが空文字の場合エラーを返す",
			userId:            validInput.userId,
			title:             validInput.title,
			githubURL:         "",
			category:          validInput.category,
			description:       validInput.description,
			tags:              validInput.tags,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid githubUrl",
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_categoryが空文字の場合エラーを返す",
			userId:            validInput.userId,
			title:             validInput.title,
			githubURL:         validInput.githubURL,
			category:          "",
			description:       validInput.description,
			tags:              validInput.tags,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid category",
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_descriptionが空文字の場合エラーを返す",
			userId:            validInput.userId,
			title:             validInput.title,
			githubURL:         validInput.githubURL,
			category:          validInput.category,
			description:       "",
			tags:              validInput.tags,
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid description",
			expectRepoNotCall: true,
		},
		{
			name:              "準正常系_tagsが空文字の場合エラーを返す",
			userId:            validInput.userId,
			title:             validInput.title,
			githubURL:         validInput.githubURL,
			category:          validInput.category,
			description:       validInput.description,
			tags:              "",
			setupMock:         func(m *repositories_blogs.MockBlogRepository) {},
			wantErr:           true,
			wantErrMsg:        "invalid tags",
			expectRepoNotCall: true,
		},
		{
			name:        "異常系_Repositoryがエラーを返す場合エラーを返す",
			userId:      validInput.userId,
			title:       validInput.title,
			githubURL:   validInput.githubURL,
			category:    validInput.category,
			description: validInput.description,
			tags:        validInput.tags,
			setupMock: func(m *repositories_blogs.MockBlogRepository) {
				m.On("CreateBlog", validInput.userId, validInput.title, validInput.githubURL, validInput.category, validInput.description, validInput.tags).Return(nil, errors.New("repository failure"))
			},
			wantErr:    true,
			wantErrMsg: "failed to create blog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositories_blogs.MockBlogRepository)
			blogService := services_blogs.NewBlogService(mockRepo)

			tt.setupMock(mockRepo)

			blog, err := blogService.CreateBlog(tt.userId, tt.title, tt.githubURL, tt.category, tt.description, tt.tags)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				assert.Nil(t, blog)
				if tt.expectRepoNotCall {
					mockRepo.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, blog)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
