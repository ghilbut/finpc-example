output endpoint {
  value = aws_lb.nlb.dns_name
}

output grpc_client_variables {
  value = <<-EOF
      GRPC_CACERT        = ${base64encode(tls_self_signed_cert.ca.cert_pem)}
      GRPC_HOST          = "${aws_lb.nlb.dns_name}"
      GRPC_HOST_OVERRIDE = "${aws_lb.alb.dns_name}"
      GRPC_INSECURE      = "false"
      GRPC_PORT          = "${aws_lb_listener.nlb_grpc.port}"
    EOF
  sensitive = true
}

output postgres_password {
  value     = random_password.rds.result
  sensitive = true
}

output postgres_session {
  value = <<EOF
aws ssm start-session \
    %{~ if var.aws_profile != "default" ~}
    --profile ${var.aws_profile} \
    %{~ endif ~}
    --region ${var.aws_region} \
    --target ${aws_instance.bastian.id} \
    --document-name AWS-StartPortForwardingSessionToRemoteHost \
    --parameters host="${aws_db_instance.this.address}",portNumber="${aws_db_instance.this.port}",localPortNumber="${local.postgres_port}"
EOF
}
