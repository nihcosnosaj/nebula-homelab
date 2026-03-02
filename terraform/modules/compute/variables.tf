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

variable "ami_id" {
  description = "The AMI ID for the instances"
  type        = string
}

variable "subnet_id" {
  description = "The Subnet ID to launch instances in"
  type        = string
}

variable "instance_profile" {
  description = "The IAM instance profile for the nodes"
  type        = string
}

variable "key_name" {
  description = "The SSH key pair name"
  type        = string
}

variable "sg_id" {
  description = "The Security Group ID for the EC2 instances"
  type        = string
}