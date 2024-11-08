default:
	go run . -in=. -out=sparkle -key=123 -header=header.txt && cat sparkle && Rscript sparkle

install:
	go install
