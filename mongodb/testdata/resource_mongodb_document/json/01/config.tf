provider "mongodb" {}

resource "mongodb_document" "main" {
  database    = "test"
  collection  = "resource.mongodb.document"
  document_id = "main"

  json = jsonencode({
    string  = "hello world"
    integer = 42
    float   = 3.14
    boolean = true
  })
}
