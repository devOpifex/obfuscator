default:
	go run . -in=test -out=test_obfuscated -key=s -license=header.txt

install:
	go install
