default:
	go run . -in=test -out=sparkle -key=123 -header=header.txt && cat sparkle/test.R; echo

install:
	go install
