# ---------------------------------------------
# IAM role
# ---------------------------------------------
# invoker_member（BFF のサービスアカウント）だけが呼べる前提。サービスの invoker IAM チェックを無効（invoker-iam-disabled）にすると
# このバインディングに関係なく誰でも呼べるようになるため、有効のまま保つ（deploy.yml の --invoker-iam-check。#165）
resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.echo_blog_app_service.name
  location = var.gcp_region
  project  = var.gcp_project_id

  role   = var.invoker_role
  member = var.invoker_member
}

resource "google_secret_manager_secret_iam_member" "cloud_run_access" {
  for_each  = local.secrets
  secret_id = google_secret_manager_secret.secrets[each.key].id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}
