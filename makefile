default:
	go run . -in=test -out=test_obfuscated -key=123 -header=header.txt && cat test_obfuscated/test.R; echo

install:
	go install
