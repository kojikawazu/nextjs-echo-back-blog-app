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
        image = "${var.gcp_region}-docker.pkg.dev/${var.gcp_project_id}/${google_artifact_registry_repository.echo_blog_app_repo.repository_id}/${var.service_name}:latest"

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

  depends_on = [
    google_artifact_registry_repository.echo_blog_app_repo,
    google_secret_manager_secret_iam_member.cloud_run_access
  ]
}


