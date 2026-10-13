# ---------------------------------------------
# GitHub Actions のデプロイ認証（Workload Identity Federation）
# ---------------------------------------------
# deploy.yml は長期有効なサービスアカウントキーを使わず、GitHub の OIDC トークンで
# 専用のデプロイ用サービスアカウントを借用する（#144）。
#
# Pool はこのリポジトリ専用に作る。principalSet（借用を許可する主体）は Pool 単位で照合されるため、
# 他リポジトリと Pool を共用すると、後から追加された条件の緩い Provider を経由して、このリポジトリの
# 別ブランチ・別ワークフローのトークンがバインディングに一致しうる。専用 Pool なら Provider はここで定義した 1 つだけになる。
#
# Cloud Run サービスは Terraform の定義が本番とずれている（#170）。-target で適用したときに
# 依存として巻き込まないよう、既存のサービス・リポジトリはリソースを参照せず変数で指定する。

resource "google_iam_workload_identity_pool" "github_echo_blog_backend" {
  workload_identity_pool_id = "github-echo-blog-backend"
  display_name              = "GitHub echo-blog backend"
  description               = "${var.github_repository} の GitHub Actions 専用"
}

# このリポジトリの main ブランチの deploy.yml が発行したトークンだけを受け入れる。
# リポジトリ名ではなく repository_id で判定する（同名で作り直されたリポジトリを通さないため）。
resource "google_iam_workload_identity_pool_provider" "github_echo_blog_backend" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_echo_blog_backend.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-echo-blog-backend"
  display_name                       = "GitHub echo-blog backend"
  description                        = "${var.github_repository} の deploy.yml（main）"

  attribute_mapping = {
    "google.subject"          = "assertion.sub"
    "attribute.repository"    = "assertion.repository"
    "attribute.repository_id" = "assertion.repository_id"
    "attribute.workflow_ref"  = "assertion.workflow_ref"
  }
  attribute_condition = "assertion.repository_id == '${var.github_repository_id}' && assertion.workflow_ref == '${var.github_repository}/.github/workflows/deploy.yml@refs/heads/main'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# デプロイ専用のサービスアカウント。権限はプロジェクト全体ではなく、操作対象のリソース単位で付与する。
resource "google_service_account" "deployer" {
  account_id   = "echo-blog-deployer"
  display_name = "GitHub Actions deployer (echo-blog-app)"
}

# GitHub Actions（専用 Pool の Provider を通ったこのリポジトリのトークン）にだけ借用を許可する
resource "google_service_account_iam_member" "deployer_workload_identity_user" {
  service_account_id = google_service_account.deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_echo_blog_backend.name}/attribute.repository_id/${var.github_repository_id}"
}

# イメージの push と、古いイメージの削除（Cleanup old images）。writer には削除権限が無い
resource "google_artifact_registry_repository_iam_member" "deployer_repo_admin" {
  location   = var.gcp_region
  repository = var.repository_id
  role       = "roles/artifactregistry.repoAdmin"
  member     = "serviceAccount:${google_service_account.deployer.email}"
}

# デプロイと --no-allow-unauthenticated（IAM ポリシーの更新）。対象はこのサービスのみ
resource "google_cloud_run_service_iam_member" "deployer_run_admin" {
  service  = var.service_name
  location = var.gcp_region
  project  = var.gcp_project_id
  role     = "roles/run.admin"
  member   = "serviceAccount:${google_service_account.deployer.email}"
}

# リビジョンの実行サービスアカウント（cloud-run-sa）として指定するための actAs。対象は cloud-run-sa のみ
resource "google_service_account_iam_member" "deployer_act_as_runtime" {
  service_account_id = google_service_account.cloud_run_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.deployer.email}"
}

# Actions secret（GCP_WORKLOAD_IDENTITY_PROVIDER / GCP_DEPLOY_SERVICE_ACCOUNT）に登録する値
output "github_workload_identity_provider" {
  value = google_iam_workload_identity_pool_provider.github_echo_blog_backend.name
}

output "github_deploy_service_account" {
  value = google_service_account.deployer.email
}
