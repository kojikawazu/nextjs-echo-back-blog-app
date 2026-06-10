# AWS App Runner → Google Cloud Run 移行レポート

## 目次

- [概要](#概要)
- [移行日](#移行日)
- [移行前後の構成](#移行前後の構成)
- [コスト効果](#コスト効果)
- [変更対象](#変更対象)
  - [1. Terraform](#1-terraform)
  - [2. GitHub Actions](#2-github-actions)
- [移行によるアプリケーションへの影響](#移行によるアプリケーションへの影響)

---

## 概要

バックエンドのホスティング基盤を AWS App Runner から Google Cloud Run へ移行した。
主な目的はランニングコストの削減である。

## 移行日

2026年3月

## 移行前後の構成

| 項目 | 移行前 | 移行後 |
|------|--------|--------|
| ホスティング | AWS App Runner | Google Cloud Run |
| IaC | Terraform (AWS provider) | Terraform (GCP provider) |
| CI/CD | GitHub Actions (AWS向け) | GitHub Actions (GCP向け) |

## コスト効果

| 項目 | 金額 |
|------|------|
| App Runner 月額コスト | 約 $12 |
| Cloud Run 月額コスト | 約 $0 |
| **月額削減額** | **$12** |
| **年間削減額** | **約 $144** |

## 変更対象

### 1. Terraform

- AWS App Runner リソース定義を削除
- GCP Cloud Run リソース定義を追加
- プロバイダーを AWS → GCP に変更

### 2. GitHub Actions

- デプロイワークフローを App Runner 向けから Cloud Run 向けに変更
- 認証方式を AWS 認証から GCP 認証（Workload Identity Federation 等）に変更
- コンテナイメージのプッシュ先を Artifact Registry に変更

## 移行によるアプリケーションへの影響

- バックエンド API の仕様変更なし
- エンドポイント URL の変更あり（CloudFlare DNS で対応）
- アプリケーションコード（Go / Echo）の変更なし
