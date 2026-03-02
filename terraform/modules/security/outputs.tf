output "key_name" {
    value = aws_key_pair.generated_key.key_name
}

output "instance_profile_name" {
    value = aws_iam_instance_profile.nebula_node_profile.name
}