module "network" {
  source = "./modules/network"
}

module "storage" {
  source = "./modules/storage"
}

module "security" {
  source = "./modules/security"
  bucket_arn = module.storage.bucket_arn 
}

module "compute" {
  source = "./modules/compute"
  ami_id = data.aws_ami.ubuntu.id
  subnet_id = module.network.public_subnet_id
  sg_id = module.network.sg_id
  instance_profile = module.security.instance_profile_name
  key_name = module.security.key_name
}

# create ansible-ready inventory of IPs for cluster config.
resource "local_file" "ansible_inventory" {
  content = <<-EOT
    [control_plane]
    ${module.compute.master_public_ip} ansible_user=ubuntu

    [workers]
    %{for ip in module.compute.worker_public_ips ~}
    ${ip} ansible_user=ubuntu ansible_ssh_private_key_file=nebula-key.pem
    %{ endfor ~}
  EOT
  filename = "${path.module}/../ansible/inventory.ini"
}





