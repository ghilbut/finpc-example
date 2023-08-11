################################################################
##
##  Github Actions - Secrets and Variables
##

resource github_actions_variable debug_actions {
  repository    = var.github_repository
  variable_name = "ACTIONS_STEP_DEBUG"
  value         = var.github_debug_actions
}

resource github_actions_variable enable_actions {
  repository    = var.github_repository
  variable_name = "ENABLE_GITHUB_ACTIONS"
  value         = var.github_enable_actions
}

resource github_actions_variable project_name {
  repository    = var.github_repository
  variable_name = "PROJECT_NAME"
  value         = var.project
}

resource github_actions_variable aws_region {
  repository    = var.github_repository
  variable_name = "AWS_REGION"
  value         = var.aws_region
}

resource github_actions_secret aws_iam_role_for_actions {
  repository      = var.github_repository
  secret_name     = "AWS_IAM_ROLE"
  plaintext_value = module.iam_github_oidc_role.arn
}

################################################################
##
##  AWS IAM  - Role for AWS ECR and ECS services
##

module iam_github_oidc_provider {
  source  = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-provider"
  version = "5.28.0"
}

module iam_github_oidc_role {
  source  = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"
  version = "5.28.0"

  name = "${var.project}-github-actions"

  subjects = [
    "${var.github_org}/${var.github_repository}:*",
  ]

  policies = {
    ECX = aws_iam_policy.github_actions.arn
  }
}

resource aws_iam_policy github_actions {
  name   = "${var.project}-github-actions"
  policy = <<-POLICY
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "GetAuthorizationToken",
          "Effect": "Allow",
          "Action": "ecr:GetAuthorizationToken",
          "Resource": "*"
        },
        {
          "Sid": "PushElasticContainerRegistryImage",
          "Effect": "Allow",
          "Action": [
            "ecr:CompleteLayerUpload",
            "ecr:BatchGetImage",
            "ecr:UploadLayerPart",
            "ecr:InitiateLayerUpload",
            "ecr:BatchCheckLayerAvailability",
            "ecr:PutImage"
          ],
          "Resource": [
            "${aws_ecr_repository.client.arn}",
            "${aws_ecr_repository.server.arn}"
          ]
        },
        {
          "Sid": "DescribeTaskDefinition",
          "Effect": "Allow",
          "Action": "ecs:DescribeTaskDefinition",
          "Resource": "*"
        },
        {
          "Sid": "RegisterTaskDefinition",
          "Effect": "Allow",
          "Action": "ecs:RegisterTaskDefinition",
          "Resource": "*"
        },
        {
           "Sid": "PassRolesInTaskDefinition",
           "Effect": "Allow",
           "Action": "iam:PassRole",
           "Resource": "${aws_iam_role.ecs_task_execution.arn}"
        },
        {
           "Sid": "DeployService",
           "Effect": "Allow",
           "Action": [
              "ecs:UpdateService",
              "ecs:DescribeServices"
           ],
           "Resource": [
              "${aws_ecs_service.client.id}",
              "${aws_ecs_service.server.id}"
           ]
        }
      ]
    }
    POLICY
}
