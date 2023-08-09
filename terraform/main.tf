terraform {
  required_version = "~> 1.5"

  backend local {
    path = "terraform.tfstate"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.11"
    }
    github = {
      source  = "integrations/github"
      version = "~> 5.32"
    }
  }
}

provider aws {
  profile = var.aws_profile
  region  = var.aws_region

  default_tags {
    tags = {
      "github.com/reop" = "https://github.com/ghilbut/finpc-example"
      "github.com/path" = "terraform/"
      project           = var.project
    }
  }
}

provider github {
  token = var.github_token
}
