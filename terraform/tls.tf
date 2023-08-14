resource tls_private_key ca {
  algorithm = "RSA"
}

resource tls_self_signed_cert ca {
  private_key_pem = tls_private_key.ca.private_key_pem

  is_ca_certificate = true

  subject {
    country             = "IN"
    province            = "Mahrashatra"
    locality            = "Mumbai"
    common_name         = "Cloud Manthan Root CA"
    organization        = "Cloud Manthan Software Solutions Pvt Ltd."
    organizational_unit = "Cloud Manthan Root Certification Auhtority"
  }

  validity_period_hours = 43800 //  1825 days or 5 years

  allowed_uses = [
    "digital_signature",
    "cert_signing",
    "crl_signing",
  ]
}

resource tls_private_key server {
  algorithm = "RSA"
}

resource tls_cert_request server {
  private_key_pem = tls_private_key.server.private_key_pem

  dns_names = [aws_lb.alb.dns_name]

  subject {
    country             = "IN"
    province            = "Mahrashatra"
    locality            = "Mumbai"
    common_name         = "Cloud Manthan Internal Development "
    organization        = "Cloud Manthan"
    organizational_unit = "Development"
  }
}

resource tls_locally_signed_cert server {
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = tls_cert_request.server.cert_request_pem
  ca_private_key_pem = tls_private_key.ca.private_key_pem

  validity_period_hours = 43800

  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "server_auth",
    "client_auth",
  ]
}
