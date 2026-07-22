package services_blogs

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockBlogService は BlogService のテスト用モック。
type MockBlogService struct {
	mock.Mock
}

// FetchBlogs は BlogService.FetchBlogs のテスト用モック。
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従う全ブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogs() ([]models.BlogData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogsByUserId は BlogService.FetchBlogsByUserId のテスト用モック。
//
// 引数:
//   - userId: 取得対象のユーザーID
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従うブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BlogData), args.Error(1)
}

// FetchBlogById は BlogService.FetchBlogById のテスト用モック。
//
// 引数:
//   - id: 取得対象のブログID
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従うブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogById(id string) (*models.BlogData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// CreateBlog は BlogService.CreateBlog のテスト用モック。
//
// 引数:
//   - userId: 投稿者のユーザーID
//   - title: ブログのタイトル
//   - githubUrl: 関連する GitHub の URL
//   - category: ブログのカテゴリ
//   - description: ブログの本文・説明
//   - tags: カンマ区切りのタグ文字列
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従う作成済みブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(userId, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// UpdateBlog は BlogService.UpdateBlog のテスト用モック。
//
// 引数:
//   - id: 更新対象のブログID
//   - title: 更新後のタイトル
//   - githubUrl: 更新後の GitHub の URL
//   - category: 更新後のカテゴリ
//   - description: 更新後の本文・説明
//   - tags: 更新後のカンマ区切りのタグ文字列
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従う更新後ブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(id, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// DeleteBlog は BlogService.DeleteBlog のテスト用モック。
//
// 引数:
//   - id: 削除対象のブログID
//
// 戻り値:
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) DeleteBlog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// FetchBlogCategories は BlogService.FetchBlogCategories のテスト用モック。
//
// 戻り値:
//   - []string: モックの戻り値設定に従うカテゴリ名の一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogCategories() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogTags は BlogService.FetchBlogTags のテスト用モック。
//
// 戻り値:
//   - []string: モックの戻り値設定に従うタグ名の一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogPopular は BlogService.FetchBlogPopular のテスト用モック。
//
// 引数:
//   - count: 取得する件数
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従う人気ブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogService) FetchBlogPopular(count int) ([]models.BlogData, error) {
	args := m.Called(count)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}
