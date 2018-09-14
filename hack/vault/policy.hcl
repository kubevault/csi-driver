# capability to create a token against the "applications" role
path "auth/token/create/applications" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
# capability to create a token against the "nginx" role
path "auth/token/create/nginx" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# capability to get role definitions (like allowed policies)
path "auth/token/roles/applications" {
  capabilities = ["read"]
}
path "auth/token/roles/nginx" {
  capabilities = ["read"]
}

# capability to list roles
path "auth/token/roles" {
  capabilities = ["read", "list"]
}