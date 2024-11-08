# obfuscator

Obfuscate R code.


```
go install github.com/sparkle-tech/obfuscator@latest
```

Flags:

- `in`: directory where to find R files to obfuscate.
- `out`: name of obfuscated file to create.
- `key`: unique key used to obfuscate code.
- `header`: header text to insert at top of obfuscated file, generally license.

Example usage

```
obfuscate -in=R -out=sparkle -key=123 -header=header.txt
```
