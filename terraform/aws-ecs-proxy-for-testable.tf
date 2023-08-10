locals {
  proxy_name    = "proxy"
  proxy_cpu     = "256"
  proxy_mem     = "512"
  postgres_port = 5432
}

resource local_file proxy {
  filename = "../proxy/nginx.conf"
  content  = templatefile("../proxy/nginx.conf.template", {
    PG_ENDPOINT = aws_db_instance.this.endpoint
  })
}

################################################################
##
##  AWS ECR
##

resource aws_ecr_repository proxy {
  name = "${var.project}-${local.proxy_name}"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${var.project}-proxy"
  }
}

################################################################
##
##  AWS ECS Fargate
##

resource aws_ecs_service proxy {
  lifecycle {
    ignore_changes = [
      task_definition,
    ]
  }

  name                               = local.proxy_name
  cluster                            = aws_ecs_cluster.this.id
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
  launch_type                        = "FARGATE"
  task_definition                    = aws_ecs_task_definition.proxy.arn
  desired_count                      = 1

  network_configuration {
    subnets         = [aws_subnet.private["a"].id]
    security_groups = [aws_security_group.proxy.id]
  }

  load_balancer {
    container_name   = local.proxy_name
    container_port   = local.postgres_port
    target_group_arn = aws_lb_target_group.nlb_proxy_postgres.id
  }

  deployment_controller {
    type = "ECS"
  }
}

resource aws_ecs_task_definition proxy {
  family                = "${var.project}-proxy"
  container_definitions = <<-JSON
    [
      {
        "name": "${local.proxy_name}",
        "image": "${aws_ecr_repository.proxy.repository_url}:latest",
        "essential": true,
        "portMappings": [
          {
            "protocol": "tcp",
            "containerPort": 80
          },
          {
            "protocol": "tcp",
            "containerPort": ${local.postgres_port}
          }
        ],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-datetime-format": "[%Y-%m-%d %H:%M:%S]",
            "awslogs-group":           "${aws_cloudwatch_log_group.this.name}",
            "awslogs-region":          "${var.aws_region}",
            "awslogs-stream-prefix":   "proxy"
          }
        }
      }
    ]
    JSON

  cpu                      = local.proxy_cpu
  memory                   = local.proxy_mem
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  tags = {
    Name = "${var.project}-proxy"
  }
}

##--------------------------------------------------------------
##  security group

resource aws_security_group proxy {
  name   = "${var.project}-ecs-proxy"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [var.cidr_block]
  }

  ingress {
    from_port   = local.postgres_port
    to_port     = local.postgres_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "sg-${var.project}-ecs-proxy"
  }
}
