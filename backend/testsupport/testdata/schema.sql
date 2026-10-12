-- IT（testcontainers）用スキーマ定義。
-- テーブル名・列は本番コード（repositories/**_impl.go・models/*.go）に一致させている。
-- 実行順は docker-entrypoint-initdb.d のファイル名昇順（schema.sql → seed.sql）。

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- updated_at を UPDATE 時に自動更新する共通トリガ関数（.claude/rules/database.md。#117）。
-- 本番（docs/supabase-migration-guide.md 付録 B）と同じ定義にする。
CREATE FUNCTION public.set_updated_at() RETURNS trigger
LANGUAGE plpgsql SET search_path = '' AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

-- ブログユーザー
CREATE TABLE blog_users (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name       varchar,
    email      varchar,
    password   varchar,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);
CREATE TRIGGER trg_blog_users_set_updated_at
    BEFORE UPDATE ON blog_users FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

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
    updated_at   timestamptz DEFAULT now() NOT NULL
);
CREATE TRIGGER trg_blogs_set_updated_at
    BEFORE UPDATE ON blogs FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ブログいいね
CREATE TABLE blog_likes (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    blog_id    uuid,
    visit_id   uuid,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);
CREATE TRIGGER trg_blog_likes_set_updated_at
    BEFORE UPDATE ON blog_likes FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ブログコメント
CREATE TABLE blog_comments (
    id         uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    blog_id    uuid,
    guest_user varchar,
    comment    varchar,
    created_at timestamptz DEFAULT now() NOT NULL
);
