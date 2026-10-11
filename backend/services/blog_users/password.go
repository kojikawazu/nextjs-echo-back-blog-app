package services_blog_users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// passwordHashCost は bcrypt のコスト。移行 SQL（pgcrypto の gen_salt('bf', 10)）と揃える。
const passwordHashCost = bcrypt.DefaultCost

// dummyPasswordHash はユーザー不在時の照合に使う bcrypt ハッシュ（どのパスワードとも一致しない）。
// ユーザー不在でも照合 1 回分の時間をかけ、応答時間からメールアドレスの登録有無を推測されにくくする。
const dummyPasswordHash = "$2a$10$0bbn8l3eXgr2GAHKWg4wgOGj4zjkBbcliA/igMQUDS.8ayp1Vrdxa"

// hashPassword はパスワードを bcrypt でハッシュ化する。
//
// 引数:
//   - password: ハッシュ化する平文のパスワード
//
// 戻り値:
//   - string: bcrypt ハッシュ
//   - error: 72 バイト超の場合は ErrPasswordTooLong、その他ハッシュ化に失敗した場合のエラー
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", ErrPasswordTooLong
		}
		return "", err
	}
	return string(hash), nil
}

// verifyPassword は bcrypt ハッシュとパスワードが一致するかを判定する。
// ハッシュが不正な形式（平文など）の場合も一致しないものとして扱う。
//
// 引数:
//   - hash: 保存されている bcrypt ハッシュ
//   - password: 照合する平文のパスワード
//
// 戻り値:
//   - bool: 一致すれば true
func verifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
