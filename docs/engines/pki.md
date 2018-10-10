vault write pki/config/urls \
    issuing_certificates="http://142.93.77.58:30001/v1/pki/ca" \
    crl_distribution_points="http://142.93.77.58:30001/v1/pki/crl"

vault write pki/roles/my-pki-role \
    allowed_domains=my-website.com \
    allow_subdomains=true \
    max_ttl=72h

vault write pki/issue/my-pki-role \
    common_name=www.my-website.com