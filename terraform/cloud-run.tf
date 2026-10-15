# ---------------------------------------------
# Cloud Run
# ---------------------------------------------
# Google Cloud Run のサービスアカウントを作成
resource "google_service_account" "cloud_run_sa" {
  account_id   = "cloud-run-sa"
  display_name = "Cloud Run Service Account"
}

# Google Cloud Run にデプロイするサービス
resource "google_cloud_run_service" "echo_blog_app_service" {
  name     = var.service_name
  location = var.gcp_region

  metadata {
    namespace = var.gcp_project_id
  }

  template {
    metadata {
      annotations = {
        # レートリミットはインスタンスのメモリでカウントするため、最大インスタンス数で実効上限を固定する（deploy.yml の --max-instances と揃える）
        "autoscaling.knative.dev/maxScale" = "3"
      }
    }

    spec {
      containers {
        # image は deploy.yml（gcloud run deploy --image <コミット SHA>）が管理するため lifecycle.ignore_changes で除外している。
        # この値はサービスの新規作成時にだけ使われ、次の deploy.yml の実行でアプリのイメージに置き換わる。
        # Artifact Registry には SHA タグしか push しないため、作成時点で必ず存在する公開のプレースホルダーを指定する（#170）
        image = "us-docker.pkg.dev/cloudrun/container/hello"

        ports {
          container_port = var.api_port
        }
        resources {
          limits = {
            cpu    = "1000m"
            memory = "512Mi"
          }
        }

        # BFF が転送する X-Client-IP を信頼する。バックエンドを BFF からのみ呼べる（invoker_member が BFF のサービスアカウントのみ）ことが前提（#146）
        env {
          name  = "TRUST_CLIENT_IP_HEADER"
          value = "true"
        }

        dynamic "env" {
          for_each = local.secrets

          content {
            name = upper(env.key)
            value_from {
              secret_key_ref {
                name = env.key
                key  = "latest"
              }
            }
          }
        }
      }
      service_account_name = google_service_account.cloud_run_sa.email
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  # リビジョンの管理の分担（#170）:
  # - deploy.yml: image（コミット SHA タグ）と、gcloud が付与するクライアント情報のアノテーション
  # - Terraform: それ以外（env・Secret 参照・リソース・maxScale・サービスアカウント）
  # 分担を決めずに両方が書くと、terraform apply がデプロイ済みのイメージを巻き戻す
  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers[0].image,
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["run.googleapis.com/client-version"],
    ]
  }

  depends_on = [
    google_artifact_registry_repository.echo_blog_app_repo,
    google_secret_manager_secret_iam_member.cloud_run_access
  ]
}


