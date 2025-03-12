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

__Caveat__

Currently cannot parse:

- `if` statements without curly braces, e.g.: `ìf(TRUE) 1 else 0`
make sure they have the surrounding curly braces: `if(TRUE) {1} else {0}`
- Functions without curly braces, e.g.: `\(x) x + 1` must be written as `\(x) {x + 1}`
- Doesn't understand expressions in curly braces (outside of the `for`, `if`, function bodies, etc.)
rewrite `tryCatch({x + 1})` to `fn <- \(){x + 1};tryCatch(fn())`
- Functions that start with a dot are __not obfuscated__, e.g.: `.onLoad`
- Files names are also obfuscated, e.g.: `foo.R` becomes `xyz.R` except files names 
starting with `__` (their content is obfuscated)
