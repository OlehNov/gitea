resource "aws_security_group" "gitea_sg" {
  name        = var.sg_name
  description = "Security Group for VPC"
  vpc_id      = var.vpc_id

  tags = {
    Name = var.sg_name
  }
}


# Ingress rules
resource "aws_security_group_rule" "ingress" {
  for_each = { for idx, rule in var.sg_ingress_rules : idx => rule }

  type              = "ingress"
  from_port         = each.value.from_port
  to_port           = each.value.to_port
  protocol          = each.value.protocol
  cidr_blocks       = each.value.cidr_blocks
  security_group_id = aws_security_group.gitea_sg.id
}

# Egress rules
resource "aws_security_group_rule" "egress" {
  for_each = { for idx, rule in var.sg_egress_rules : idx => rule }

  type              = "egress"
  from_port         = each.value.from_port
  to_port           = each.value.to_port
  protocol          = each.value.protocol
  cidr_blocks       = each.value.cidr_blocks
  security_group_id = aws_security_group.gitea_sg.id
}
