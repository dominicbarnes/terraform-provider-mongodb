# terraform-provider-mongodb

> A [Terraform][terraform] provider for [MongoDB][mongodb]. The primary use-case
> for this particular provider is "seeding" a database for test/stage/dev
> environments.

## provider "mongodb"

| Attribute | Type | Description |
| --- | --- | --- |
| hostname | string | The hostname of MongoDB server to connect to (default 'localhost') |
| password | string | The password of the user connecting to MongoDB |
| port | number | The port number of the MongoDB server to connect to (default 27017) |
| username | string | The name of the user connecting to MongoDB |


### resource "mongodb_document"

| Attribute | Type | Description |
| --- | --- | --- |
| collection | string | The collection for the MongoDB document |
| database | string | The database for the MongoDB document |
| document_id | string | The _id value for the MongoDB document |
| json | string | The raw JSON of the MongoDB document without _id |


[terraform]: https://www.terraform.io/
[mongodb]: https://www.mongodb.com/