resource "aws_instance" "control_plane" {
  ami           = var.ami_id
  instance_type = var.control_plane_type

  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = "0.03"
    }
  }

  subnet_id              = var.subnet_id
  vpc_security_group_ids = [var.sg_id]
  key_name               = var.key_name

  tags = {
    Name = "nebula-control-plane"
    Project = var.project_name
    Role = "control-plane"
  }
}

resource "aws_instance" "worker_nodes" {
  count = var.worker_count
  depends_on = [ aws_instance.control_plane ]
  ami = var.ami_id
  instance_type = var.worker_type

  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = "0.08"
    }
  }

  subnet_id              = var.subnet_id
  vpc_security_group_ids = [var.sg_id]
  key_name               = var.key_name

  tags = {
    Name = "nebula-worker-${count.index}"
    Project = var.project_name
    Role = "worker"
  }
}