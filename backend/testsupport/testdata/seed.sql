-- IT（testcontainers）用の決定的シードデータ。
-- 固定 UUID を使い、testsupport/testdb.go の fixture 定数および
-- TEST_* 環境変数と一致させている（本番ダンプの実データは流用しない）。

INSERT INTO blog_users (id, name, email, password) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Test User', 'test@example.com', 'test-password');

-- 2件用意する（FetchBlogs / FetchBlogsByUserId が len>=2 を期待するため）。両方とも同一ユーザー所有。
INSERT INTO blogs (id, blog_user_id, title, description, github_url, category, tags, likes, comment_cnt) VALUES
    ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111',
     'Test Blog Title', 'Test blog description', 'https://github.com/test/repo', 'Tech', 'go,testing', 0, 0),
    ('22222222-2222-2222-2222-222222222223', '11111111-1111-1111-1111-111111111111',
     'Second Blog Title', 'Second blog description', 'https://github.com/test/repo2', 'Go', 'echo,api', 0, 0);

INSERT INTO blog_likes (id, blog_id, visit_id) VALUES
    ('33333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222',
     '44444444-4444-4444-4444-444444444444');

INSERT INTO blog_comments (id, blog_id, guest_user, comment) VALUES
    ('55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', 'Guest User', 'Nice post!');
