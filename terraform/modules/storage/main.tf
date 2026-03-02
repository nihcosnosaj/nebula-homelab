resource "random_id" "suffix" {
  byte_length = 4
}

resource "aws_s3_bucket" "nebula_data" {
  bucket = "nebula-homelab-storage-${random_id.suffix.hex}"
}

output "bucket_arn" { value = aws_s3_bucket.nebula_data.arn }
output "bucket_name" { value = aws_s3_bucket.nebula_data.bucket }