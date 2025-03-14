# obfuscator

Obfuscate R code.

> [!WARNING]
> This tool only obfuscates code and does not encrypt it.
> Obfuscation is not a security measure and should not be relied upon to protect 
> sensitive code or intellectual property.
> The obfuscated code can potentially be reverse-engineered.
> Do not use this tool with the expectation that it makes your code secure to share.

```
go install github.com/sparkle-tech/obfuscator@latest
```

## Usage

```
$> obfuscator -h
Usage of obfuscator:
  -decipher string
        String to decypher
-in string
    Directory of R files to obfuscate
-key string
    Key to obfuscate
-license string
    License to prepend to every obfuscated file, e.g.: license
-out string
    Directory where to write the obfuscated files
-protect string
    Comma separated protected tokens, e.g.: foo,bar
```

__Example__

```
obfuscate -in=R -out=obfuscated -key=123 -license=license.txt
```

## Caveats

Currently known issues and usage:

- `if` statements without curly braces, e.g.: `ìf(TRUE) 1 else 0`
make sure they have the surrounding curly braces: `if(TRUE) {1} else {0}`
- Functions without curly braces, e.g.: `\(x) x + 1` must be written as `\(x) {x + 1}`
- Doesn't understand expressions in curly braces (outside of the `for`, `if`, function bodies, etc.)
rewrite `tryCatch({x + 1})` to `fn <- \(){x + 1};tryCatch(fn())`
- Functions that start with a dot are __not obfuscated__, e.g.: `.onLoad`
- Files names are also obfuscated, e.g.: `foo.R` becomes `xyz.R` __except__ files names 
starting with `__` (their content is obfuscated)
- Only `.R` files are obfuscated
- arguments `do.call` is __not obfuscated__, avoid `do.call`.
