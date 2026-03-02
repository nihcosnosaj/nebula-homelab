resource "aws_kms_key" "vault_unseal" {
  description             = "KMS key for Vault auto-unseal in Nebula"
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_iam_role" "nebula_node_role" {
  name = "nebula-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy" "s3_access" {
  name = "nebula-s3-access"
  role = aws_iam_role.nebula_node_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["s3:PutObject", "s3:GetObject", "s3:ListBucket"]
        Effect   = "Allow"
        Resource = [var.bucket_arn, "${var.bucket_arn}/*"] # Ideally scoped to your specific bucket ARN
      }
    ]
  })
}

resource "aws_iam_instance_profile" "nebula_node_profile" {
  name = "nebula-node-profile"
  role = aws_iam_role.nebula_node_role.name
}

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
  filename             = "${path.root}/nebula-key.pem"
  file_permission      = "0400"
  directory_permission = "0700"
}