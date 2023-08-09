################################################################
##
##  Github Actions - Secrets and Variables
##

resource github_actions_variable enable_actions {
  repository    = var.github_repository
  variable_name = "ENABLE_GITHUB_ACTIONS"
  value         = "true"
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
  plaintext_value = ""
}
