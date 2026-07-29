# nextjs-echo-back-blog-app

Go + Echo フレームワークで構築されたブログWebアプリケーションのバックエンドAPI

## Rules

明示的な指示がなくても、`.claude/rules/` 内のルールを常に守ってください。

| ファイル | スコープ | 内容 |
|---------|---------|------|
| shortcuts.md | 全体 | 指示ショートカット（PR出して、PR承認しました 等） |
| workflow.md | 全体 | 開発フロー（ブランチ運用・テスト必須） |
| quality-gate.md | 全体 | 品質ゲート（セルフレビュー・設計/実装レビュー） |
| documentation.md | 全体 | ドキュメント更新ルール |
| git.md | 全体 | GitHub Flow・ブランチ命名・push 禁止物 |
| github-issue.md | 全体 | GitHub issue 運用（ブランチと対で起票・open/close で進捗管理・`Closes #N` 自動クローズ・サブ issue） |
| testing.md | 全体 | テスト分類・原則・Go テストツール/配置 |
| coding-standards.md | 全体 | コーディング規約（Go 1.22・go mod・golangci-lint/gofmt） |
| error-handling.md | 全体 | エラーハンドリング方針（バリデーション・HTTP ステータス・統一レスポンス） |
| security.md | 全体 | セキュリティ方針（認証・CORS・インジェクション対策・シークレット管理） |
| api.md | backend/** | Go + Echo API 設計・レイヤー分離（Handler/Service/Repository） |
| database.md | backend/** | DB 方針（監査列の DB 側自動化・timestamptz 統一・スキーマ同期） |
