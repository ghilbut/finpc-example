################################################################
##
##  AWS ELB
##

##--------------------------------------------------------------
##  external network load balancer

resource aws_lb nlb {
  name                             = "nlb-${var.project}"
  enable_cross_zone_load_balancing = true
  internal                         = false
  load_balancer_type               = "network"
  subnets                          = [ for k, v in aws_subnet.public : v.id ]

  tags = {
    Name = "nlb-${var.project}"
  }
}

resource aws_lb_listener nlb_http {
  load_balancer_arn = aws_lb.nlb.arn
  port              = aws_lb_target_group.nlb_http.port
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.nlb_http.arn
  }

  tags = {
    Name = "${var.project}-nlb-http"
  }
}

resource aws_lb_target_group_attachment nlb_http {
  port             = aws_lb_listener.alb_http.port
  target_group_arn = aws_lb_target_group.nlb_http.arn
  target_id        = aws_lb.alb.id
}

resource aws_lb_target_group nlb_http {
  name        = "nlb-${var.project}-http"
  port        = aws_lb_listener.alb_http.port
  protocol    = "TCP"
  target_type = "alb"
  vpc_id      = aws_vpc.this.id

  health_check {
    interval = 10
    matcher  = "200,404"
  }
}

resource aws_lb_listener nlb_grpc {
  load_balancer_arn = aws_lb.nlb.arn
  port              = aws_lb_target_group.nlb_grpc.port
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.nlb_grpc.arn
  }

  tags = {
    Name = "${var.project}-nlb-grpc"
  }
}

resource aws_lb_target_group_attachment nlb_grpc {
  port             = aws_lb_listener.alb_grpc.port
  target_group_arn = aws_lb_target_group.nlb_grpc.arn
  target_id        = aws_lb.alb.id
}

resource aws_lb_target_group nlb_grpc {
  name        = "nlb-${var.project}-grpc"
  port        = aws_lb_listener.alb_grpc.port
  protocol    = "TCP"
  target_type = "alb"
  vpc_id      = aws_vpc.this.id

  health_check {
    interval = 10
    matcher  = "200,404"
    port     = aws_lb_listener.alb_http.port
  }
}

##--------------------------------------------------------------
##  internal application load balancer

resource aws_lb alb {
  name               = "alb-${var.project}"
  internal           = true
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = [ for k, v in aws_subnet.private : v.id ]

  tags = {
    Name = "alb-${var.project}",
  }
}

resource aws_lb_listener alb_http {
  load_balancer_arn = aws_lb.alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found (AWS Application Load Balancer)"
      status_code  = "404"
    }
  }

  tags = {
    Name = "${var.project}-alb-http"
  }
}

resource aws_lb_listener alb_grpc {
  load_balancer_arn = aws_lb.alb.arn
  port              = local.grpc_port
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = aws_acm_certificate.grpc.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alb_grpc.arn
  }

  tags = {
    Name = "${var.project}-alb-grpc"
  }
}

resource aws_security_group alb {
  name   = "${var.project}-alb"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = local.grpc_port
    to_port     = local.grpc_port
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
    Name = "sg-${var.project}-alb",
  }
}

################################################################
##
##  AWS Certificate Manager
##

resource aws_acm_certificate grpc {
  private_key      = tls_private_key.server.private_key_pem
  certificate_body = tls_locally_signed_cert.server.cert_pem
}
