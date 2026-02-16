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
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
    enable_dns_support = true

    tags = {
        Name = "nebula-vpc"
        Project = "nebula"
    }
}

resource "aws_internet_gateway" "igw" {
    vpc_id = aws_vpc.nebula_vpc.id
}

resource "aws_subnet" "public_subnet" {
    vpc_id = aws_vpc.nebula_vpc.id
    cidr_block = "10.0.1.0/24"
    map_public_ip_on_launch = true
    availability_zone = "us-west-1a"

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
    subnet_id = aws_subnet.public_subnet.id
    route_table_id = aws_route_table.public_rt.id
}

resource "aws_security_group" "nebula_sg" {
    name = "nebula-sg"
    description = "Allow SSH and interal k8s traffic"
    vpc_id = aws_vpc.nebula_vpc.id

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"] # restrict to just my IP when in prod
    }
    
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

# Spot Instance req (aiming for 16GB RAM)
resource "aws_instance" "worker_node" {
    ami = data.aws_ami.ubuntu.id
    instance_type = "t4g.xlarge"

    # request spot within the instance resource
    instance_market_options {
        market_type = "spot"
        spot_options {
            max_price = "0.08"
            spot_instance_type = "one-time"
        }
    }

    subnet_id = aws_subnet.public_subnet.id
    vpc_security_group_ids = [aws_security_group.nebula_sg.id]
    key_name = aws_key_pair.generated_key.key_name

    tags = {
        Name = "nebula-node"
        Project = "nebula"
    }
}

# get secure SSH key
resource "tls_private_key" "nebula_key" {
    algorithm = "ED25519"
}

# register key with AWS
resource "aws_key_pair" "generated_key" {
    key_name = "nebula-key"
    public_key = tls_private_key.nebula_key.public_key_openssh
}

# save key to local disk 
resource "local_sensitive_file" "private_key" {
    content = tls_private_key.nebula_key.private_key_openssh
    filename = "${path.module}/nebula-key.pem"
    file_permission = "0400"
    directory_permission = "0700"
}