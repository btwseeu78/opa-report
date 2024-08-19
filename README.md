## The codes uses local kubeconfig to generate violation report pod based violation will be converted to deployment name and only one instance will be there.
### Running the code
`` go get . ``
`` go mod tidy ``
`` go run main.go get`

### building the binary
The codes used cobra to create acmmand line tool.

``go build .``
#### it will generate atyifact called opa-report(exculatvle)
./opa-report get
