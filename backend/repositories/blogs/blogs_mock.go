package repositories_blogs

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockBlogRepository は BlogRepository のテスト用モック。
type MockBlogRepository struct {
	mock.Mock
}

// FetchBlogs は BlogRepository.FetchBlogs のモック実装。
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従う全ブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogs() ([]models.BlogData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogsByUserId は BlogRepository.FetchBlogsByUserId のモック実装。
//
// 引数:
//   - userId: 取得対象のユーザーID
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従うブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogById は BlogRepository.FetchBlogById のモック実装。
//
// 引数:
//   - id: 取得対象のブログID
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従うブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogById(id string) (*models.BlogData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateBlog は BlogRepository.CreateBlog のモック実装。
//
// 引数:
//   - userId: 作成者のユーザーID
//   - title: ブログのタイトル
//   - githubUrl: 関連する GitHub の URL
//   - category: ブログのカテゴリ
//   - description: ブログの説明
//   - tags: ブログのタグ
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従う作成後のブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(userId, title, githubUrl, category, description, tags)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateBlog は BlogRepository.UpdateBlog のモック実装。
//
// 引数:
//   - id: 更新対象のブログID
//   - title: 更新後のタイトル
//   - githubUrl: 更新後の GitHub の URL
//   - category: 更新後のカテゴリ
//   - description: 更新後の説明
//   - tags: 更新後のタグ
//
// 戻り値:
//   - *models.BlogData: モックの戻り値設定に従う更新後のブログデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(id, title, githubUrl, category, description, tags)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// DeleteBlog は BlogRepository.DeleteBlog のモック実装。
//
// 引数:
//   - id: 削除対象のブログID
//
// 戻り値:
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) DeleteBlog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// FetchBlogCategories は BlogRepository.FetchBlogCategories のモック実装。
//
// 戻り値:
//   - []string: モックの戻り値設定に従うカテゴリ一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogCategories() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogTags は BlogRepository.FetchBlogTags のモック実装。
//
// 戻り値:
//   - []string: モックの戻り値設定に従うタグ一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogPopular は BlogRepository.FetchBlogPopular のモック実装。
//
// 引数:
//   - count: 取得する件数
//
// 戻り値:
//   - []models.BlogData: モックの戻り値設定に従う人気ブログデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogRepository) FetchBlogPopular(count int) ([]models.BlogData, error) {
	args := m.Called(count)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}
