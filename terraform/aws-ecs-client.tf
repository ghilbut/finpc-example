locals {
  client_name = "client"
  client_cpu  = "256"
  client_mem  = "512"
  client_port = 3000
}

################################################################
##
##  AWS ECR
##

resource aws_ecr_repository client {
  name = "${var.project}-${local.client_name}"

  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${var.project}-client"
  }
}

################################################################
##
##  AWS ECS Fargate
##

resource aws_ecs_service client {
  lifecycle {
    ignore_changes = [
      task_definition,
    ]
  }

  name                               = local.client_name
  cluster                            = aws_ecs_cluster.this.id
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
  launch_type                        = "FARGATE"
  task_definition                    = aws_ecs_task_definition.client.arn
  desired_count                      = 1

  network_configuration {
    subnets          = [aws_subnet.private["a"].id]
    security_groups  = [aws_security_group.client.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.client.id
    container_name   = local.client_name
    container_port   = local.client_port
  }

  deployment_controller {
    type = "ECS"
  }
}

resource aws_ecs_task_definition client {
  family                = "${var.project}-client"
  container_definitions = <<JSON
    [
      {
        "name": "${local.client_name}",
        "image": "${aws_ecr_repository.client.repository_url}:latest",
        "essential": true,
        "environment": [
          {
            "name": "GRPC_HOST",
            "value": "${aws_lb.alb.dns_name}"
          },
          {
            "name": "GRPC_PORT",
            "value": "${local.grpc_port}"
          },
          {
            "name": "SENTRY_ENVIRONMENT",
            "value": "ecs"
          }
        ],
        "secrets": [
          {
            "name": "GRPC_CACERT",
            "valueFrom": "${aws_secretsmanager_secret.grpc_cacert.arn}"
          },
          {
            "name": "SENTRY_DSN",
            "valueFrom": "${aws_secretsmanager_secret.sentry_client_dsn.arn}"
          }
        ],
        "portMappings": [
          {
            "protocol": "tcp",
            "containerPort": ${local.client_port}
          }
        ],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-create-group":    "true",
            "awslogs-datetime-format": "[%Y-%m-%d %H:%M:%S]",
            "awslogs-group":           "${aws_cloudwatch_log_group.this.name}",
            "awslogs-region":          "${var.aws_region}",
            "awslogs-stream-prefix":   "ecs"
          }
        }
      }
    ]
    JSON

  cpu                      = local.client_cpu
  memory                   = local.client_mem
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  tags = {
    Name = "${var.project}-client"
  }
}

##--------------------------------------------------------------
##  security group

resource aws_security_group client {
  name   = "${var.project}-ecs-client"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port       = local.client_port
    to_port         = local.client_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "sg-${var.project}-ecs-client"
  }
}

################################################################
##
##  AWS Secrets Manager
##

resource aws_secretsmanager_secret grpc_cacert {
  name = "${var.project}-grpc-cacert"
  recovery_window_in_days = 0
}

resource aws_secretsmanager_secret_version grpc_cacert {
  secret_id     = aws_secretsmanager_secret.grpc_cacert.id
  secret_string = base64encode(tls_self_signed_cert.ca.cert_pem)
}

resource aws_secretsmanager_secret sentry_client_dsn {
  name = "${var.project}-sentry-client-dsn"
  recovery_window_in_days = 0
}

resource aws_secretsmanager_secret_version sentry_client_dsn {
  secret_id     = aws_secretsmanager_secret.sentry_client_dsn.id
  secret_string = var.sentry_client_dsn
}

################################################################
##
##  AWS ELB
##

resource aws_lb_listener_rule client {
  listener_arn = aws_lb_listener.alb_http.arn
  priority     = local.alb_client_priority

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.client.arn
  }

  condition {
    path_pattern {
      values = ["/*"]
    }
  }
}

resource aws_lb_target_group client {
  name        = "alb-${var.project}-client"
  port        = local.client_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.this.id

  deregistration_delay = 10

  health_check {
    interval = 10
    path = "/"
  }
}
