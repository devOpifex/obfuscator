# obfuscator

Obfuscate R code.


```
go install github.com/sparkle-tech/obfuscator@latest
```

__Flags__

- `in`: directory where to find R files to obfuscate.
- `out`: name of obfuscated file to create.
- `key`: unique key used to obfuscate code.
- `header`: header text to insert at top of obfuscated file, generally license (optional).

```
$> obfuscator -h
Usage of obfuscator:
  -header string
        Header to append to obfuscated code, e.g.: license
  -in string
        Directory of R files to obfuscate
  -key string
        Key to obfuscate
  -out string
        Output name of file to write obfuscated code
```

__Example usage__

```
obfuscate -in=R -out=sparkle -key=123 -header=header.txt
```
