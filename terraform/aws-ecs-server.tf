locals {
  server_name = "server"
  server_cpu  = "256"
  server_mem  = "512"
  rest_port   = 8080
  grpc_port   = 9095
}

################################################################
##
##  AWS ECR
##

resource aws_ecr_repository server {
  name = "${var.project}-${local.server_name}"

  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${var.project}-server"
  }
}

################################################################
##
##  AWS ECS Fargate
##

resource aws_ecs_service server {
  lifecycle {
    ignore_changes = [
      task_definition,
    ]
  }

  name                               = local.server_name
  cluster                            = aws_ecs_cluster.this.id
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
  launch_type                        = "FARGATE"
  task_definition                    = aws_ecs_task_definition.server.arn
  desired_count                      = 1

  network_configuration {
    subnets         = [aws_subnet.private["a"].id]
    security_groups = [aws_security_group.server.id]
  }

  load_balancer {
    container_name   = local.server_name
    container_port   = local.rest_port
    target_group_arn = aws_lb_target_group.alb_rest.id
  }

  load_balancer {
    container_name   = local.server_name
    container_port   = local.grpc_port
    target_group_arn = aws_lb_target_group.alb_grpc.id
  }

  deployment_controller {
    type = "ECS"
  }
}

resource aws_ecs_task_definition server {
  family                = "${var.project}-server"
  container_definitions = <<-JSON
    [
      {
        "name": "${local.server_name}",
        "image": "${aws_ecr_repository.server.repository_url}:latest",
        "essential": true,
        "environment": [
          {
            "name": "PG_HOST",
            "value": "${aws_db_instance.this.address}"
          },
          {
            "name": "PG_PORT",
            "value": "${aws_db_instance.this.port}"
          },
          {
            "name": "PG_USER",
            "value": "postgres"
          },
          {
            "name": "PG_DATABASE",
            "value": "postgres"
          },
          {
            "name": "PG_SSLMODE",
            "value": "require"
          },
          {
            "name": "SENTRY_ENVIRONMET",
            "value": "ecs"
          }
        ],
        "secrets": [
          {
            "name": "PG_PASSWORD",
            "valueFrom": "${aws_secretsmanager_secret.postgres_password.arn}"
          },
          {
            "name": "SENTRY_DSN",
            "valueFrom": "${aws_secretsmanager_secret.sentry_server_dsn.arn}"
          }
        ],
        "portMappings": [
          {
            "protocol": "tcp",
            "containerPort": ${local.rest_port}
          },
          {
            "protocol": "tcp",
            "containerPort": ${local.grpc_port}
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

  cpu                      = local.server_cpu
  memory                   = local.server_mem
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  tags = {
    Name = "${var.project}-server"
  }
}

##--------------------------------------------------------------
##  security group

resource aws_security_group server {
  name        = "${var.project}-ecs-server"
  description = "Allow all traffic"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port       = local.rest_port
    to_port         = local.rest_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    from_port       = local.grpc_port
    to_port         = local.grpc_port
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
    Name = "sg-${var.project}-ecs-server"
  }
}

################################################################
##
##  AWS Secrets Manager
##

resource aws_secretsmanager_secret postgres_password {
  name = "${var.project}-postgres-password"
  recovery_window_in_days = 0
}

resource aws_secretsmanager_secret_version postgres_password {
  secret_id     = aws_secretsmanager_secret.postgres_password.id
  secret_string = random_password.rds.result
}

resource aws_secretsmanager_secret sentry_server_dsn {
  name = "${var.project}-sentry-server-dsn"
  recovery_window_in_days = 0
}

resource aws_secretsmanager_secret_version sentry_server_dsn {
  secret_id     = aws_secretsmanager_secret.sentry_server_dsn.id
  secret_string = var.sentry_server_dsn
}

################################################################
##
##  AWS ELB
##

resource aws_lb_listener_rule alb_rest {
  listener_arn = aws_lb_listener.alb_http.arn
  priority     = local.alb_rest_priority

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alb_rest.arn
  }

  condition {
    path_pattern {
      values = ["/v1/*"]
    }
  }

  tags = {
    Name = "${var.project}-alb-rest-server"
  }
}

resource aws_lb_target_group alb_rest {
  name        = "alb-${var.project}-rest-server"
  port        = local.rest_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.this.id

  deregistration_delay = 10

  health_check {
    interval = 10
    path = "/healthz"
  }
}

resource aws_lb_target_group alb_grpc {
  name             = "alb-${var.project}-grpc-server"
  port             = local.grpc_port
  protocol         = "HTTP"
  protocol_version = "GRPC"
  target_type      = "ip"
  vpc_id           = aws_vpc.this.id

  deregistration_delay = 10

  health_check {
    interval = 10
    matcher  = "12"
    path     = "/"
  }
}
