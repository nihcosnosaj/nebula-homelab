output "public_subnet_id" {
    value = aws_subnet.public_subnet.id
}

output "sg_id" {
    value = aws_security_group.nebula_sg.id
}

output "vpc_id" {
    value = aws_vpc.nebula_vpc.id
}