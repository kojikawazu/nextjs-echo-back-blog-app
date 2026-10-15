# AWS App Runner → Google Cloud Run 移行レポート

## 目次

- [概要](#概要)
- [移行日](#移行日)
- [移行前後の構成](#移行前後の構成)
- [コスト効果](#コスト効果)
- [変更対象](#変更対象)
  - [1. Terraform](#1-terraform)
  - [2. GitHub Actions](#2-github-actions)
  - [3. Workload Identity Federation への移行](#3-workload-identity-federation-への移行)
  - [4. Terraform と本番の差分（drift）の解消](#4-terraform-と本番の差分drift-の解消)
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
- 認証方式を AWS 認証から GCP 認証に変更（移行時点ではサービスアカウントキー。Workload Identity Federation への移行は 3. を参照）
- コンテナイメージのプッシュ先を Artifact Registry に変更

### 3. Workload Identity Federation への移行

移行時点の `deploy.yml` は、サービスアカウントキー（Actions secret `GCP_SERVICE_ACCOUNT_KEY`）で認証していた。長期有効な鍵は漏洩すると失効させるまで悪用できるため、Workload Identity Federation（WIF）でデプロイ専用サービスアカウントを借用する方式に移行した（#144）。定義は `terraform/github-actions-wif.tf`。

| 項目 | 移行前 | 移行後 |
|------|--------|--------|
| 認証 | サービスアカウントキー（JSON） | WIF（GitHub の OIDC トークン） |
| サービスアカウント | 他リポジトリと共用の SA（プロジェクト全体に `run.admin` / `secretmanager.admin` 等） | 専用の `echo-blog-deployer`（リソース単位の最小権限） |
| 古いイメージの削除 | 削除権限が無く毎回失敗（イメージが蓄積） | `repoAdmin` + `--delete-tags` で削除 |

#### 切り替え手順

元に戻せる操作を先に、戻せない操作（鍵の削除）を最後に行う。

1. **Terraform の適用（追加のみ）**。Cloud Run サービスの定義が本番とずれているため（#170）、対象リソースに限定して適用する。plan が `7 to add, 0 to change, 0 to destroy` であることを確認してから適用する。

   ```bash
   cd terraform
   terraform apply \
     -target=google_iam_workload_identity_pool.github_echo_blog_backend \
     -target=google_iam_workload_identity_pool_provider.github_echo_blog_backend \
     -target=google_service_account.deployer \
     -target=google_service_account_iam_member.deployer_workload_identity_user \
     -target=google_artifact_registry_repository_iam_member.deployer_repo_admin \
     -target=google_cloud_run_service_iam_member.deployer_run_admin \
     -target=google_service_account_iam_member.deployer_act_as_runtime
   ```

2. **Actions secret の登録**（値は Terraform の output から渡し、画面に出さない）

   ```bash
   gh secret set GCP_WORKLOAD_IDENTITY_PROVIDER --body "$(terraform output -raw github_workload_identity_provider)"
   gh secret set GCP_DEPLOY_SERVICE_ACCOUNT --body "$(terraform output -raw github_deploy_service_account)"
   ```

3. **`deploy.yml` の変更をマージ**し、`main` への push でデプロイが成功することを確認する。監査ログ（Cloud Run の `ReplaceService`）の `principalEmail` が `echo-blog-deployer` で、`serviceAccountKeyName` が空（鍵を使っていない）であることを確認する。失敗した場合は PR を revert すれば鍵の認証に戻る（鍵は手順 4 まで残っている）。
4. **鍵の廃止（戻せない）**。Actions secret `GCP_SERVICE_ACCOUNT_KEY` を削除し、旧サービスアカウントのユーザー管理キーを削除する。旧サービスアカウント自体は他リポジトリが WIF で使っているため削除しない。

   ```bash
   gh secret delete GCP_SERVICE_ACCOUNT_KEY
   gcloud iam service-accounts keys list --iam-account <旧SA> --managed-by=user
   gcloud iam service-accounts keys delete <キーID> --iam-account <旧SA>
   ```

### 4. Terraform と本番の差分（drift）の解消

Cloud Run のリビジョンは `deploy.yml`（`gcloud run deploy`）が更新するため、Terraform の定義が追従せず、`terraform plan` が `1 to change` になっていた（#170）。この状態で `terraform apply` を全体にかけると、デプロイ済みのイメージを巻き戻す。

| 差分 | 原因 | 対応 |
|------|------|------|
| イメージ（Terraform は `:latest`、本番はコミット SHA タグ） | `deploy.yml` が SHA タグで更新する（`:latest` は push していない） | `deploy.yml` の管理とし、`lifecycle.ignore_changes` で除外。Terraform の値は新規作成時のプレースホルダー（`us-docker.pkg.dev/cloudrun/container/hello`）にした |
| アノテーション `run.googleapis.com/client-name` / `client-version` | gcloud が付与する | `lifecycle.ignore_changes` で除外 |
| env の並び | 本番の並びが手作業で設定されたもので、Terraform の定義（`TRUST_CLIENT_IP_HEADER` → Secret 参照をキー順）と異なる。`env` は順序付きのリストとして比較される | Terraform の定義を正とし、適用して本番を揃える |
| env `DUMMY` | 本番だけにある。アプリケーションは参照していない | 適用で削除する |

#### 解消手順

1. plan で差分が env（並びと `DUMMY`）だけであることを確認する。イメージ・アノテーションが差分に出ないこと（`ignore_changes` が効いていること）を確認する。

   ```bash
   cd terraform
   terraform plan -target=google_cloud_run_service.echo_blog_app_service
   ```

2. Cloud Run サービスだけに限定して適用する。新しいリビジョンが作られ、イメージは現在の SHA タグのまま、env が Terraform の定義どおりになる。

   ```bash
   terraform apply -target=google_cloud_run_service.echo_blog_app_service
   ```

3. 全体の `terraform plan` が `No changes` になることを確認する。以降は `-target` を付けずに適用できる。

## 移行によるアプリケーションへの影響

- バックエンド API の仕様変更なし
- エンドポイント URL の変更あり（CloudFlare DNS で対応）
- アプリケーションコード（Go / Echo）の変更なし
