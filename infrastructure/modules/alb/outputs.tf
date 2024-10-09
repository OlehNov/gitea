output "alb_dns_name" {
  value = aws_lb.this.dns_name
}

# Outputs for Target Groups
output "django_target_group_arn" {
  value = aws_lb_target_group.golang_target_group.arn
}
