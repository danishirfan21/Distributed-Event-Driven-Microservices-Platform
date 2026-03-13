provider "aws" {
  region = var.region
}

resource "aws_db_instance" "postgres" {
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "15.3"
  instance_class       = "db.t3.micro"
  db_name              = "ecommerce"
  username             = "postgres"
  password             = var.db_password
  skip_final_snapshot  = true
}

resource "aws_eks_cluster" "main" {
  name     = "dist-ecommerce"
  role_arn = aws_iam_role.eks_cluster.arn

  vpc_config {
    subnet_ids = var.subnet_ids
  }
}

# NATS is typically self-hosted on EKS or used as a managed service like NGS.
# Here we provide a placeholder for a self-hosted NATS on EKS.
# Alternatively, use a managed NATS provider.

resource "helm_release" "nats" {
  name       = "nats"
  repository = "https://nats-io.github.io/k8s/helm/charts/"
  chart      = "nats"
  namespace  = "messaging"
  create_namespace = true

  set {
    name  = "nats.jetstream.enabled"
    value = "true"
  }
}
