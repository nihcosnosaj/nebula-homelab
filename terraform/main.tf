# fetch AMI dynamically.
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical's official AWS ID

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-arm64-server-*"]
  }
}

resource "aws_vpc" "nebula_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name    = "nebula-vpc"
    Project = var.project_name
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.nebula_vpc.id
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.nebula_vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-west-1a"

  tags = {
    Name = "nebula-public"
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.nebula_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "public_asoc" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_security_group" "nebula_sg" {
  name        = "nebula-sg"
  description = "Allow SSH and interal k8s traffic"
  vpc_id      = aws_vpc.nebula_vpc.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # restrict to just my IP when in prod
  }

  ingress {
    description = "k3s API server"
    from_port = 6443
    to_port = 6443
    protocol = "tcp"
    cidr_blocks = [aws_vpc.nebula_vpc.cidr_block]
  }

  ingress {
    description = "internal cluster traffic"
    from_port = 0
    to_port = 0
    protocol = "-1"
    self = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_role" "control_plane_role" {
  name = "nebula-control-plane-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}


# control-plane node resource request
resource "aws_instance" "control_plane" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = var.control_plane_type

  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = "0.03"
    }
  }

  subnet_id              = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.nebula_sg.id]
  key_name               = aws_key_pair.generated_key.key_name

  tags = {
    Name = "nebula-control-plane"
    Project = var.project_name
    Role = "control-plane"
  }
}

# worker nodes (2 instances)
resource "aws_instance" "worker_nodes" {
  count = var.worker_count
  depends_on = [ aws_instance.control_plane ]
  ami = data.aws_ami.ubuntu.id
  instance_type = var.worker_type

  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = "0.08"
    }
  }

  subnet_id              = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.nebula_sg.id]
  key_name               = aws_key_pair.generated_key.key_name

  tags = {
    Name = "nebula-worker-${count.index}"
    Project = var.project_name
    Role = "worker"
  }
}

# get secure SSH key
resource "tls_private_key" "nebula_key" {
  algorithm = "ED25519"
}

# register key with AWS
resource "aws_key_pair" "generated_key" {
  key_name   = "nebula-key"
  public_key = tls_private_key.nebula_key.public_key_openssh
}

# save key to local disk 
resource "local_sensitive_file" "private_key" {
  content              = tls_private_key.nebula_key.private_key_openssh
  filename             = "${path.module}/nebula-key.pem"
  file_permission      = "0400"
  directory_permission = "0700"
}

# create ansible-ready inventory of IPs for cluster config.
resource "local_file" "ansible_inventory" {
  content = <<-EOT
    [control_plane]
    ${aws_instance.control_plane.public_ip} ansible_user=ubuntu

    [workers]
    %{for ip in aws_instance.worker_nodes[*].public_ip ~}
    ${ip} ansible_user=ubuntu ansible_ssh_private_key_file=nebula-key.pem
    %{ endfor ~}
  EOT
  filename = "${path.module}/../ansible/inventory.ini"
}