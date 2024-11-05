default:
	go run . -in=. -out=sparkle -key=123 -header=header.txt && Rscript sparkle

install:
	go install
