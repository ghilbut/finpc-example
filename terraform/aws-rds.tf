################################################################
##
##  AWS RDS - PostgreSQL
##

resource aws_db_instance this {
  allocated_storage      = 20
  db_subnet_group_name   = aws_db_subnet_group.this.name
  engine                 = "postgres"
  engine_version         = "15.3"
  identifier             = var.project
  instance_class         = "db.t4g.micro"
  network_type           = "IPV4"
  storage_type           = "gp3"
  username               = "postgres"
  password               = random_password.rds.result
  skip_final_snapshot    = true
  vpc_security_group_ids = [aws_security_group.rds.id]
}

resource aws_db_subnet_group this {
  name       = var.project
  subnet_ids = values(aws_subnet.private).*.id

  tags = {
    Name = "${var.project}-db-subnet-group"
  }
}

resource aws_security_group rds {
  name   = "${var.project}-postgres"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
    security_groups = [
      aws_security_group.bastian.id,
    ]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = -1
    security_groups = [
      aws_security_group.bastian.id,
    ]
  }

  ingress {
    from_port   = 5432
    to_port     = 5432
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
    Name = "sg-${var.project}-rds-postgres",
  }
}

##--------------------------------------------------------------
##  password

resource random_password rds {
  length = 12

  min_lower   = 2
  min_numeric = 2
  min_special = 2
  min_upper   = 2

  override_special = "!#$%&*()-_=+[]{}<>:?"
}
