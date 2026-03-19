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
  type = string
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