package repositories_blogs

import "backend/models"

// BlogRepository はブログのデータアクセスを担うリポジトリインターフェース。
type BlogRepository interface {
	// FetchBlogs は全ブログデータを取得する。
	//
	// 戻り値:
	//   - []models.BlogData: 取得した全ブログデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogs() ([]models.BlogData, error)
	// FetchBlogsByUserId は指定されたユーザーIDに一致するブログデータを取得する。
	//
	// 引数:
	//   - userId: 取得対象のユーザーID
	//
	// 戻り値:
	//   - []models.BlogData: 取得したブログデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	// FetchBlogById は指定されたIDに一致するブログデータを取得する。
	//
	// 引数:
	//   - id: 取得対象のブログID
	//
	// 戻り値:
	//   - *models.BlogData: 取得したブログデータ
	//   - error: 取得に失敗した場合のエラー
	FetchBlogById(id string) (*models.BlogData, error)

	// CreateBlog はブログデータを作成する。
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
	//   - *models.BlogData: 作成したブログデータ
	//   - error: 作成に失敗した場合のエラー
	CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// UpdateBlog はブログデータを更新する。
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
	//   - *models.BlogData: 更新後のブログデータ
	//   - error: 更新に失敗した場合のエラー
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// DeleteBlog はブログデータを削除する。
	//
	// 引数:
	//   - id: 削除対象のブログID
	//
	// 戻り値:
	//   - error: 削除に失敗した場合のエラー
	DeleteBlog(id string) error

	// FetchBlogCategories はブログカテゴリ一覧を取得する。
	//
	// 戻り値:
	//   - []string: 取得したカテゴリ一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogCategories() ([]string, error)
	// FetchBlogTags はブログタグ一覧を取得する。
	//
	// 戻り値:
	//   - []string: 取得したタグ一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogTags() ([]string, error)
	// FetchBlogPopular は人気のあるブログを取得する。
	//
	// 引数:
	//   - count: 取得する件数
	//
	// 戻り値:
	//   - []models.BlogData: 取得した人気ブログデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogPopular(count int) ([]models.BlogData, error)
}

// BlogRepositoryImpl は BlogRepository の実装。
type BlogRepositoryImpl struct{}

// NewBlogRepository は BlogRepository を実装した BlogRepositoryImpl のポインタを生成する。
//
// 戻り値:
//   - BlogRepository: 生成したリポジトリ実装
func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImpl{}
}
