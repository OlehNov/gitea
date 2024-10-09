module "vpc" {
  source = "/home/oleh/VIsual_studio/gitea/infrastructure/modules/networks/vpc"
}

resource "aws_eks_cluster" "gitea_cluster" {
  name     = "gitea_cluster"
  role_arn = aws_iam_role.eks_role.arn



  vpc_config {
    subnet_ids = module.vpc.private_subnet_ids
#     vpc_id = module.vpc.vpc_id
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster_policy,
    aws_iam_role_policy_attachment.eks_vpc_resource_controller,
  ]
}


resource "aws_iam_role" "eks_role" {
  name = "test_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRoleWithWebIdentity"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })

  tags = {
    tag-key = "gitea_teg"
  }
}

# Привязка необходимых политик к роли IAM для EKS
resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_role.name
}

resource "aws_iam_role_policy_attachment" "eks_vpc_resource_controller" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.eks_role.name
}



