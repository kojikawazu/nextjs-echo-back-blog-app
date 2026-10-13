# ---------------------------------------------
# Variables
# ---------------------------------------------
variable "gcp_project_id" {
  type = string
  sensitive = true 
}

variable "gcp_region" {
  type = string
}

variable "repository_id" {
  type = string
}

variable "invoker_role" {
  type = string
}

variable "invoker_member" {
  type        = string
  description = "バックエンドを呼び出せる主体。BFF（フロントエンドの Cloud Run）のサービスアカウントを指定する（例: serviceAccount:cloud-run-sa@<フロントのプロジェクト>.iam.gserviceaccount.com）。allUsers にしない（#146）"
}

variable "service_name" {
  type = string
}

variable "api_port" {
  type = number
}

variable "allowed_origins" {
  type = string
  sensitive = true 
}

variable "env_word" {
  type = string
  sensitive = true 
}

variable "jwt_secret_key" {
  type = string
  sensitive = true 
}

variable "supabase_url" {
  type = string
  sensitive = true 
}

variable "tmp_allowed_origins" {
  type = string
  sensitive = true 
}

variable "github_repository" {
  type        = string
  default     = "kojikawazu/nextjs-echo-back-blog-app"
  description = "デプロイを許可する GitHub リポジトリ（owner/name）。Provider の条件で workflow_ref と照合する"
}

variable "github_repository_id" {
  type        = string
  default     = "870653888"
  description = "デプロイを許可する GitHub リポジトリの数値 ID（gh api repos/<owner>/<name> --jq .id）。同名で作り直されたリポジトリを通さないため、名前ではなく ID で判定する"
}
