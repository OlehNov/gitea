module "vpc" {
  source = "/home/oleh/VIsual_studio/gitea/infrastructure/modules/networks/vpc"
}

# ALB Security Group
resource "aws_security_group" "this" {
  name        = "my-alb-sg"
  description = "Allow inbound HTTP traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Application Load Balancer
resource "aws_lb" "this" {
  name                       = "gitea-alb"
  internal                   = false
  load_balancer_type         = "application"
  security_groups            = [aws_security_group.this.id]
  subnets                    = module.vpc.public_subnet_ids
#   enable_deletion_protection = var.enable_deletion_protection
  enable_deletion_protection = false
  enable_http2               = true

  tags = {
    Name = var.alb_name
  }
}


# Target Group for Django
resource "aws_lb_target_group" "golang_target_group" {
  name     = "django-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = module.vpc.vpc_id

  health_check {
    interval            = 30
    path                = "/"
    protocol            = "HTTP"
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
}

# Listener for HTTP traffic
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.golang_target_group.arn
  }
}

# Rule to forward traffic to Golang based on path
resource "aws_lb_listener_rule" "golang_rule" {
  listener_arn = aws_lb_listener.http.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.golang_target_group.arn
  }

  condition {
    path_pattern {
      values = ["/golang*"]
    }
  }
}
