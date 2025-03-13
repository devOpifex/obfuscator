default:
	go run . -in=test -out=test_obfuscated -key=123 -license=header.txt

install:
	go install
