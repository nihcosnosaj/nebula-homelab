variable "aws_region" {
    type = string
    default = "us-west-1"
}

variable "project_name" {
    type = string
    default = "nebula"
}

variable "control_plane_type" {
    type = string
    default = "t4g.medium"
}

variable "worker_type" {
  type = string
  default = "t4g.large"
}

variable "worker_count" {
    type = number
    default = 2
}