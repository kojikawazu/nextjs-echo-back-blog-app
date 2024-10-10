# ---------------------------------------------
# Secret Manager - 環境変数
# ---------------------------------------------
resource "aws_secretsmanager_secret" "echo_env" {
  name = "echo-blog-app-env-secret"
}

resource "aws_secretsmanager_secret_version" "echo_env" {
  secret_id = aws_secretsmanager_secret.echo_env.id
  secret_string = jsonencode({
    ALLOWED_ORIGINS = "${var.cors_address}",
    PORT            = "${var.api_port}",
    SUPABASE_URL    = "${var.supabase_url}",
    JWT_SECRET_KEY  = "${var.jwt_secret_key}",
    ENV             = "${var.env_word}",
  })
}

data "aws_secretsmanager_secret" "echo_env" {
  name = aws_secretsmanager_secret.echo_env.name
}

data "aws_secretsmanager_secret_version" "echo_env" {
  secret_id = data.aws_secretsmanager_secret.echo_env.id
}
