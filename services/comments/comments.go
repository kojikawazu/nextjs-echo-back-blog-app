package services_comments

import (
	"backend/models"
	"errors"
	"log"
)

// 指定されたブログIDに一致するコメントデータを取得する
func (s *CommentServiceImpl) FetchCommentsByBlogId(blogId string) ([]models.CommentData, error) {
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

// コメントデータを新規作成する
func (s *CommentServiceImpl) CreateComment(blogId, guestUser, comment string) (*models.CommentData, error) {
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
