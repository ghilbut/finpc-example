################################################################
##
##  AWS EC2 - Bastian server
##

data aws_ami bastian {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "architecture"
    values = ["arm64"]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-2023*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource aws_instance bastian {
  ami               = data.aws_ami.bastian.id
  availability_zone = aws_subnet.private["a"].availability_zone

  credit_specification {
    cpu_credits = "unlimited"
  }
  ebs_optimized        = true
  iam_instance_profile = aws_iam_instance_profile.bastian.id
  instance_type        = "t4g.nano"
  root_block_device {
    volume_size = 8
    volume_type = "gp3"
  }

  subnet_id              = aws_subnet.private["a"].id
  vpc_security_group_ids = [aws_security_group.bastian.id]

  tags = {
    Name = "${var.project}-bastian"
  }
  volume_tags = {
    Name = "${var.project}-bastian"
  }
}

resource aws_security_group bastian {
  name   = "${var.project}-bastian"
  vpc_id = aws_vpc.this.id

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = local.postgres_port
    to_port     = local.postgres_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = local.grpc_port
    to_port     = local.grpc_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "sg-${var.project}-bastian",
  }
}

resource aws_iam_instance_profile bastian {
  name = "${var.project}-bastian-instance-profile"

  role = aws_iam_role.bastian.id

  tags = {
    Name = "${var.project}-bastian-instance-profile"
  }
}

resource aws_iam_role bastian {
  name               = "${var.project}-bastian"
  assume_role_policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        }
      ]
    }
    EOF

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
  ]
}
