resource "aws_vpc" "nebula_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name    = "nebula-vpc"
    Project = var.project_name
  }
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

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.nebula_vpc.id
  
  tags = {
    Name    = "nebula-igw"
    Project = var.project_name
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.nebula_vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-west-1a"

  tags = {
    Name    = "nebula-public"
    Project = var.project_name
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.nebula_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "public_assoc" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}