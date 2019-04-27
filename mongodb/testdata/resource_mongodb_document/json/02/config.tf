provider "mongodb" {}

resource "mongodb_document" "main" {
  database    = "test"
  collection  = "resource.mongodb.document"
  document_id = "main"

  json = jsonencode({
    string  = "hello world!"
    integer = 1
    float   = 6.28
    boolean = false
  })
}
