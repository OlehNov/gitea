variable "region" {
  description = "region"
  default     = "eu-central-1"
}

variable "alb_name" {
  description = "The name of the Application Load Balancer"
  type        = string
}

variable "subnets" {
  description = "List of subnet IDs for the ALB"
  type        = list(string)
}

variable "vpc_id" {
  description = "The VPC ID where the ALB and Target Group are located"
  type        = string
}

variable "target_group_name" {
  description = "The name of the Target Group"
  type        = string
}

# variable "ec2_instance_ids" {
#   description = "List of EC2 instance IDs to attach to the Target Group"
#   type        = list(string)
# }

variable "enable_deletion_protection" {
  description = "Whether to enable deletion protection on the ALB"
  type        = bool
  default     = false
}

variable "golang_path" {
  description = "Path to route traffic to Golang"
  type        = string
  default     = "/golang*"
}
