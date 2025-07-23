module github.com/youngprinnce/product-microservice/tests

go 1.24.0

replace github.com/youngprinnce/product-microservice => ../

require (
	github.com/stretchr/testify v1.10.0
	github.com/youngprinnce/product-microservice v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
