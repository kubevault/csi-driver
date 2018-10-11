# capability to create a token against the "nginx" role
path "auth/token/create/nginx" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "auth/token/roles/nginx" {
  capabilities = ["read"]
}

# capability to list roles
path "auth/token/roles" {
  capabilities = ["read", "list"]
}

# capability of get secret
path "kv/*" {
  capabilities = ["read"]
}

# capability to get aws cred
path "aws/*" {
  capabilities = ["read"]
}

# capability to get pki cred
path "pki/*" {
  capabilities = ["read", "create", "update", "delete"]
}