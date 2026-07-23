-- IT（testcontainers）用スキーマ定義。
-- テーブル名・列は本番コード（repositories/**_impl.go・models/*.go）に一致させている。
-- 実行順は docker-entrypoint-initdb.d のファイル名昇順（schema.sql → seed.sql）。

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ブログユーザー
CREATE TABLE blog_users (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name       varchar,
    email      varchar,
    password   varchar,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamp   DEFAULT now()
);

-- ブログ記事（likes / comment_cnt はカウンタ列。本番同様 bigint）
CREATE TABLE blogs (
    id           uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    blog_user_id uuid,
    title        varchar,
    description  varchar,
    github_url   varchar,
    category     varchar,
    tags         varchar,
    likes        bigint DEFAULT 0,
    comment_cnt  bigint DEFAULT 0,
    created_at   timestamptz DEFAULT now() NOT NULL,
    updated_at   timestamp   DEFAULT now()
);

-- ブログいいね
CREATE TABLE blog_likes (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    blog_id    uuid,
    visit_id   uuid,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamp   DEFAULT now()
);

-- ブログコメント
CREATE TABLE blog_comments (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    blog_id    uuid,
    guest_user varchar,
    comment    varchar,
    created_at timestamptz DEFAULT now() NOT NULL
);
