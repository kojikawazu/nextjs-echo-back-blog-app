package services_blogs

import (
	"backend/logger"
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

// isValidUUID は値が UUID 形式かを判定する（DB に不正な形式の値を渡さないため）。
//
// 引数:
//   - value: 判定する値
//
// 戻り値:
//   - bool: UUID 形式であれば true
func isValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

// FetchBlogs は全ブログデータを取得する。
//
// 戻り値:
//   - []models.BlogData: 取得した全ブログデータの一覧
//   - error: 取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogs() ([]models.BlogData, error) {
	return s.BlogRepository.FetchBlogs()
}

// FetchBlogsByUserId は指定されたユーザーIDに一致するブログデータを取得する。
//
// 引数:
//   - userId: 取得対象のユーザーID
//
// 戻り値:
//   - []models.BlogData: 取得したブログデータの一覧
//   - error: userId が空・UUID 形式でない場合は ErrInvalidUserID、その他取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogsByUserId start...")

	// バリデーション
	if !isValidUUID(userId) {
		logger.ErrorLog.Printf("invalid userId: %s", userId)
		return nil, ErrInvalidUserID
	}
	logger.InfoLog.Println("Valid userId")

	// リポジトリを呼び出してブログデータを取得
	blogs, err := s.BlogRepository.FetchBlogsByUserId(userId)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", err)
		return nil, fmt.Errorf("fetch blogs by userId: %w", err)
	}

	logger.InfoLog.Printf("Fetched blogs successfully: %v", blogs)
	return blogs, nil
}

// FetchBlogById は指定されたIDに一致するブログデータを取得する。
//
// 引数:
//   - id: 取得対象のブログID
//
// 戻り値:
//   - *models.BlogData: 見つかったブログデータ
//   - error: id が空・UUID 形式でない場合は ErrInvalidID、不在の場合は ErrBlogNotFound、その他取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogById(id string) (*models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogById start...")

	// バリデーション
	if !isValidUUID(id) {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return nil, ErrInvalidID
	}
	logger.InfoLog.Println("Valid id")

	// リポジトリを呼び出してブログデータを取得
	blog, err := s.BlogRepository.FetchBlogById(id)
	if err != nil {
		if errors.Is(err, repositories_blogs.ErrBlogNotFound) {
			return nil, ErrBlogNotFound
		}
		logger.ErrorLog.Printf("Failed to fetch blog: %v", err)
		return nil, fmt.Errorf("fetch blog: %w", err)
	}

	logger.InfoLog.Printf("Fetched blog successfully: %v", blog)
	return blog, nil
}

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
//   - error: 入力が不正な場合は ErrInvalidXxx、その他作成に失敗した場合のエラー
func (s *BlogServiceImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("CreateBlog start...")

	// バリデーション
	if userId == "" {
		logger.ErrorLog.Printf("invalid userId: %s", userId)
		return nil, ErrInvalidUserID
	}
	if title == "" {
		logger.ErrorLog.Printf("invalid title: %s", title)
		return nil, ErrInvalidTitle
	}
	if githubUrl == "" {
		logger.ErrorLog.Printf("invalid githubUrl: %s", githubUrl)
		return nil, ErrInvalidGithubURL
	}
	if category == "" {
		logger.ErrorLog.Printf("invalid category: %s", category)
		return nil, ErrInvalidCategory
	}
	if description == "" {
		logger.ErrorLog.Printf("invalid description: %s", description)
		return nil, ErrInvalidDescription
	}
	if tags == "" {
		logger.ErrorLog.Printf("invalid tags: %s", tags)
		return nil, ErrInvalidTags
	}
	logger.InfoLog.Println("Valid input")

	// リポジトリを呼び出してブログデータを作成
	blog, err := s.BlogRepository.CreateBlog(userId, title, githubUrl, category, description, tags)
	if err != nil {
		logger.ErrorLog.Printf("Failed to create blog: %v", err)
		return nil, fmt.Errorf("create blog: %w", err)
	}

	logger.InfoLog.Printf("Created blog successfully: %v", blog)
	return blog, nil
}

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
//   - error: 入力が不正な場合は ErrInvalidXxx、不在の場合は ErrBlogNotFound、その他更新に失敗した場合のエラー
func (s *BlogServiceImpl) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("UpdateBlog start...")

	// バリデーション
	if !isValidUUID(id) {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return nil, ErrInvalidID
	}
	if title == "" {
		logger.ErrorLog.Printf("invalid title: %s", title)
		return nil, ErrInvalidTitle
	}
	if githubUrl == "" {
		logger.ErrorLog.Printf("invalid githubUrl: %s", githubUrl)
		return nil, ErrInvalidGithubURL
	}
	if category == "" {
		logger.ErrorLog.Printf("invalid category: %s", category)
		return nil, ErrInvalidCategory
	}
	if description == "" {
		logger.ErrorLog.Printf("invalid description: %s", description)
		return nil, ErrInvalidDescription
	}
	if tags == "" {
		logger.ErrorLog.Printf("invalid tags: %s", tags)
		return nil, ErrInvalidTags
	}
	logger.InfoLog.Println("Valid input")

	// リポジトリを呼び出してブログデータを更新
	blog, err := s.BlogRepository.UpdateBlog(id, title, githubUrl, category, description, tags)
	if err != nil {
		if errors.Is(err, repositories_blogs.ErrBlogNotFound) {
			return nil, ErrBlogNotFound
		}
		logger.ErrorLog.Printf("Failed to update blog: %v", err)
		return nil, fmt.Errorf("update blog: %w", err)
	}

	logger.InfoLog.Printf("Updated blog successfully: %v", blog)
	return blog, nil
}

// DeleteBlog は指定されたIDに一致するブログデータを削除する。
//
// 引数:
//   - id: 削除対象のブログID
//
// 戻り値:
//   - error: id が空・UUID 形式でない場合は ErrInvalidID、その他削除に失敗した場合のエラー
func (s *BlogServiceImpl) DeleteBlog(id string) error {
	logger.InfoLog.Printf("DeleteBlog start...")

	// バリデーション
	if !isValidUUID(id) {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return ErrInvalidID
	}
	logger.InfoLog.Println("Valid id")

	// リポジトリを呼び出してブログデータを削除
	err := s.BlogRepository.DeleteBlog(id)
	if err != nil {
		logger.ErrorLog.Printf("Failed to delete blog: %v", err)
		return fmt.Errorf("delete blog: %w", err)
	}

	logger.InfoLog.Println("Deleted blog successfully")
	return nil
}

// FetchBlogCategories はブログカテゴリを取得する。
//
// 戻り値:
//   - []string: 取得したカテゴリ名の一覧
//   - error: 取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogCategories() ([]string, error) {
	return s.BlogRepository.FetchBlogCategories()
}

// FetchBlogTags はブログタグを取得する。
//
// 戻り値:
//   - []string: 重複排除・ソート済みのタグ名一覧
//   - error: 取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogTags() ([]string, error) {
	tagsList, err := s.BlogRepository.FetchBlogTags()
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blog tags: %v", err)
		return nil, err
	}

	// カンマ区切りのタグを分割し、余分な空白を削除
	uniqueTags := make(map[string]struct{})
	for _, tags := range tagsList {
		splitTags := strings.Split(tags, ",")
		for _, tag := range splitTags {
			cleanedTag := strings.TrimSpace(tag)
			if cleanedTag != "" {
				uniqueTags[cleanedTag] = struct{}{}
			}
		}
	}

	// ユニークなタグをリスト化し、ソート
	result := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		result = append(result, tag)
	}
	sort.Strings(result)

	// ソートしたタグを返す
	return result, nil
}

// FetchBlogPopular は人気のあるブログを取得する。
//
// 引数:
//   - count: 取得する件数
//
// 戻り値:
//   - []models.BlogData: 取得した人気ブログデータの一覧
//   - error: count が 0 以下の場合は ErrInvalidCount、その他取得に失敗した場合のエラー
func (s *BlogServiceImpl) FetchBlogPopular(count int) ([]models.BlogData, error) {

	// バリデーション
	if count <= 0 {
		logger.ErrorLog.Printf("invalid count: %d", count)
		return nil, ErrInvalidCount
	}
	logger.InfoLog.Println("Valid count")

	// リポジトリを呼び出して人気のあるブログを取得
	blogs, err := s.BlogRepository.FetchBlogPopular(count)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch popular blogs: %v", err)
		return nil, fmt.Errorf("fetch popular blogs: %w", err)
	}

	logger.InfoLog.Printf("Fetched popular blogs successfully: %v", blogs)
	return blogs, nil
}
