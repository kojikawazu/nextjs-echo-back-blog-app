# ---------------------------------------------
# Variables
# ---------------------------------------------
variable "project" {
  type = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "vpc_address" {
  type = string
}

variable "default_route" {
  type = string
}

variable "public_1a_address" {
  type = string
}

variable "public_1c_address" {
  type = string
}

variable "private_1a_address" {
  type = string
}

variable "private_1c_address" {
  type = string
}

variable "api_port" {
  type = number
}

variable "cors_address" {
  type = string
}

variable "supabase_url" {
  type = string
}

variable "jwt_secret_key" {
  type = string
}

variable "env_word" {
  type = string
}
