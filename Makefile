all:
	go-bindata -debug=false -pkg="template" -o="template/bindata.go" -ignore="/\." template/code/
	go build
