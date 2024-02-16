# initialize connection to the salesforce system
provider "salesforce" {
  api_host      = "https://xyz.my.salesforce.com"
  api_version   = "v59.0"
  auth_host     = "https://login.salesforce.com/services/oauth2/token"
  client_id     = "idisjisjisjsfjs"
  client_secret = "sajfaspojfapopsa"
  grant_type    = "password"
  username      = "admin"
  password      = "password123"
}
