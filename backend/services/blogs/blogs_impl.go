package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
)

// BlogService はブログに関する処理を提供するサービスインターフェース。
type BlogService interface {
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
	//   - *models.BlogData: 見つかったブログデータ
	//   - error: 取得に失敗した場合のエラー
	FetchBlogById(id string) (*models.BlogData, error)

	// CreateBlog はブログデータを作成する。
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
	//   - *models.BlogData: 作成されたブログデータ
	//   - error: 作成に失敗した場合のエラー
	CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// UpdateBlog は指定されたIDに一致するブログデータを更新する。
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
	//   - *models.BlogData: 更新後のブログデータ
	//   - error: 更新に失敗した場合のエラー
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// DeleteBlog は指定されたIDに一致するブログデータを削除する。
	//
	// 引数:
	//   - id: 削除対象のブログID
	//
	// 戻り値:
	//   - error: 削除に失敗した場合のエラー
	DeleteBlog(id string) error

	// FetchBlogCategories はブログカテゴリを取得する。
	//
	// 戻り値:
	//   - []string: 取得したカテゴリ名の一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogCategories() ([]string, error)
	// FetchBlogTags はブログタグを取得する。
	//
	// 戻り値:
	//   - []string: 重複排除・ソート済みのタグ名一覧
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

// BlogServiceImpl は BlogService の実装。
type BlogServiceImpl struct {
	BlogRepository repositories_blogs.BlogRepository
}

// NewBlogService は BlogService インターフェースを実装した BlogServiceImpl を生成する。
//
// 引数:
//   - blogRepository: ブログデータへのアクセスを担うリポジトリ
//
// 戻り値:
//   - BlogService: 生成されたブログサービス
func NewBlogService(
	blogRepository repositories_blogs.BlogRepository,
) BlogService {
	return &BlogServiceImpl{
		BlogRepository: blogRepository,
	}
}
