output "vpc_id" {
  description = "ID VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "ID public subnet"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "ID private subnet"
  value       = aws_subnet.private[*].id
}
