provider "mongodb" {}

resource "mongodb_document" "main" {
  database   = "test"
  collection = "resource.mongodb.document"
}
