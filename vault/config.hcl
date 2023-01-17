plugin_directory = "./vault/plugins"
api_addr         = "http://127.0.0.1:8200"

storage "inmem" {}

ui = true
listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true"
}