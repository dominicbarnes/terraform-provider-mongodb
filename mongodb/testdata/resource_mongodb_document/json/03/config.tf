provider "mongodb" {}

resource "mongodb_document" "main" {
  database    = "test"
  collection  = "resource.mongodb.document"
  document_id = "main"

  json = jsonencode({
    list = ["a", "b", "c"],
    map  = { a = "A" }
  })
}
