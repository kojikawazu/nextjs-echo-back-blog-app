# secrets.tf

locals {
  secrets = {
    allowed_origins     = var.allowed_origins
    env_word            = var.env_word
    jwt_secret_key      = var.jwt_secret_key
    supabase_url        = var.supabase_url
    tmp_allowed_origins = var.tmp_allowed_origins
    gcp_project_id      = var.gcp_project_id
  }
}

resource "google_secret_manager_secret" "secrets" {
  for_each  = local.secrets
  secret_id = each.key

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_versions" {
  for_each    = local.secrets
  secret      = google_secret_manager_secret.secrets[each.key].id
  secret_data = each.value
}
