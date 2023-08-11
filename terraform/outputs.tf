output endpoint {
  value = aws_lb.nlb.dns_name
}

output grpc_cacert {
  value = base64encode(tls_self_signed_cert.ca.cert_pem)
  sensitive = true
}

output grpc_session {
  value = <<EOF
aws ssm start-session \
    %{~ if var.aws_profile != "default" ~}
    --profile ${var.aws_profile} \
    %{~ endif ~}
    --region ${var.aws_region} \
    --target ${aws_instance.bastian.id} \
    --document-name AWS-StartPortForwardingSessionToRemoteHost \
    --parameters host="${aws_lb.alb.dns_name}",portNumber="${local.grpc_port}",localPortNumber="${local.grpc_port}"
EOF
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
