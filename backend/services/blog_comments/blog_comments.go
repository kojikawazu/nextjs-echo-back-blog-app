package services_blog_comments

import (
	"backend/models"
	"errors"
	"log"
)

// FetchCommentsByBlogId は指定されたブログIDに一致するコメントデータを取得する。
//
// 引数:
//   - blogId: 取得対象のブログID
//
// 戻り値:
//   - []models.BlogCommentsData: 取得したコメントデータの一覧
//   - error: 取得に失敗した場合のエラー
func (s *CommentServiceImpl) FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error) {
	log.Printf("FetchCommentsByBlogId start...")

	// バリデーション
	if blogId == "" {
		log.Printf("invalid blogId: %s", blogId)
		return nil, errors.New("invalid blogId")
	}
	log.Println("Valid blogId")

	// リポジトリを呼び出してブログデータを取得
	comments, err := s.CommentRepository.FetchCommentsByBlogId(blogId)
	if err != nil {
		log.Printf("Failed to fetch comments: %v", err)
		return nil, errors.New("comments not found")
	}

	log.Printf("Fetched comments successfully: %v", comments)
	return comments, nil
}

// CreateComment はコメントデータを新規作成する。
//
// 引数:
//   - blogId: コメント対象のブログID
//   - guestUser: コメントを投稿するゲストユーザー名
//   - comment: コメント本文
//
// 戻り値:
//   - *models.BlogCommentsData: 作成されたコメントデータ
//   - error: 作成に失敗した場合のエラー
func (s *CommentServiceImpl) CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error) {
	log.Printf("CreateComment start...")

	// バリデーション
	if blogId == "" {
		log.Printf("invalid blogId: %s", blogId)
		return nil, errors.New("invalid blogId")
	}
	if guestUser == "" {
		log.Printf("invalid guestUser: %s", guestUser)
		return nil, errors.New("invalid guestUser")
	}
	if comment == "" {
		log.Printf("invalid comment: %s", comment)
		return nil, errors.New("invalid comment")
	}
	log.Println("Valid blogId, guestUser and comment")

	// リポジトリを呼び出してコメントデータを作成
	newComment, err := s.CommentRepository.CreateComment(blogId, guestUser, comment)
	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		return nil, errors.New("failed to create comment")
	}

	log.Printf("Created comment successfully: %v", newComment)
	return newComment, nil
}
