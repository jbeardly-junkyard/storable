all:
	go-bindata -debug=false -pkg="mongogen" -o="bindata.go" -ignore="/\." templates/
	go build
