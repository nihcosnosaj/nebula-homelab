output "ssh_command" {
  value       = "ssh -i nebula-key.pem ubuntu@${aws_instance.control_plane.public_ip}"
  description = "To log in to a node"
}