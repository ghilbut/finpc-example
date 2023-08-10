locals {
  alb_client_priority = 200
  alb_grpc_priority   = 150
  alb_rest_priority   = 100
}


variable project {
  type    = string
  default = "finpc"
}

variable aws_profile {
  type    = string
  default = "default"
}

variable aws_region {
  type    = string
  default = "ap-northeast-2"
}

variable cidr_block {
  type    = string
  default = "10.0.0.0/16"
}

variable github_token {
  type      = string
  sensitive = true
}

variable github_debug_actions {
  type    = bool
  default = false
}

variable github_enable_actions {
  type    = bool
  default = true
}

variable github_org {
  type    = string
  default = "ghilbut"
}

variable github_repository {
  type    = string
  default = "finpc-example"
}
