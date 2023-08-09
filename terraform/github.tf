################################################################
##
##  Github Actions - Secrets and Variables
##

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

resource github_actions_variable aws_access_key {
  repository    = var.github_repository
  variable_name = "AWS_ACCESS_KEY"
  value         = var.aws_access_key
}

resource github_actions_secret aws_secret_key {
  repository      = var.github_repository
  secret_name     = "AWS_SECRET_KEY"
  plaintext_value = var.aws_secret_key
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

module iam_github_oidc_role {
  source  = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"
  version = "5.28.0"

  name = "${var.project}-github-actions"

  subjects = [
    "terraform-aws-modules/terraform-aws-iam:*"
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
          "Effect": "Allow",
          "Action": [
            "ecr:PutImage"
          ],
          "Resource": [
            "${aws_ecr_repository.client.arn}",
            "${aws_ecr_repository.proxy.arn}",
            "${aws_ecr_repository.server.arn}"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "ecs:DescribeTaskDefinition",
            "ecs:RegisterTaskDefinition"
          ],
          "Resource": [
            "${aws_ecs_task_definition.client.arn}",
            "${aws_ecs_task_definition.proxy.arn}",
            "${aws_ecs_task_definition.server.arn}"
          ]
        }
      ]
    }
    POLICY
}
