################################################################
##
##  AWS ECS Fargate
##

resource aws_ecs_cluster this {
  name               = var.project

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = var.project
  }
}

resource aws_ecs_cluster_capacity_providers this {
  cluster_name = aws_ecs_cluster.this.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

################################################################
##
##  AWS IAM
##

data aws_iam_policy_document assume_role_policy {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

##--------------------------------------------------------------
##  task execution

resource aws_iam_role ecs_task_execution {
  name               = "${var.project}-ecs-task-execution"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
  ]

  inline_policy {
    name = "pull-container-images"
    policy = <<-POLICY
      {
        "Version": "2012-10-17",
        "Statement": [
          {
            "Effect": "Allow",
            "Action": [
              "ecr:GetAuthorizationToken",
              "ecr:GetDownloadUrlForLayer",
              "ecr:BatchCheckLayerAvailability",
              "ecr:BatchGetImage"
            ],
            "Resource": "*"
          }
        ]
      }
      POLICY
  }

  inline_policy {
    name = "write-logs"
    policy = <<-POLICY
      {
        "Version": "2012-10-17",
        "Statement": [
          {
            "Effect": "Allow",
            "Action": [
              "logs:CreateLogGroup",
              "logs:CreateLogStream",
              "logs:PutLogEvents",
              "logs:DescribeLogStreams"
            ],
            "Resource": "arn:aws:logs:*:*:*"
          }
        ]
      }
      POLICY
  }

  inline_policy {
    name = "read-secrets"
    policy = <<-EOF
      {
        "Version": "2012-10-17",
        "Statement": [
          {
            "Effect": "Allow",
            "Action": [
              "kms:Decrypt",
              "secretsmanager:GetSecretValue",
              "ssm:GetParameters"
            ],
            "Resource": [
              "arn:aws:secretsmanager:${var.aws_region}:*:secret:*",
              "arn:aws:ssm:${var.aws_region}:*:parameter/*",
              "arn:aws:kms:${var.aws_region}:*:key/*"
            ]
          }
        ]
      }
      EOF
  }
}
