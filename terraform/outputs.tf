output "ssh_command" {
  value       = "ssh -i nebula-key.pem ubuntu@${module.compute.master_public_ip}"
  description = "To log in to a node"
}

output "master_ip" {
  description = "Public IP of the control plane"
  value = module.compute.master_public_ip
}

output "worker_ips" {
  description = "Publid IPs of the worker nodes"
  value = module.compute.worker_public_ips
}