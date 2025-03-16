default:
	rm -r test_obfuscated
	rm -r deobfuscated
	mkdir test_obfuscated
	mkdir deobfuscated
	go run . -in=test -out=test_obfuscated -key=secret -license=header.txt
	go run . -in=test_obfuscated -out=deobfuscated -key=secret -deobfuscate=true

install:
	go install
