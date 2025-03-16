# obfuscator

A tool to obfuscate R code by renaming variables, functions, and file names while preserving functionality.

> [!WARNING]
> This tool only obfuscates code and does not encrypt it.
> Obfuscation is not a security measure and should not be relied upon to protect 
> sensitive code or intellectual property.
> The obfuscated code can potentially be reverse-engineered.
> Do not use this tool with the expectation that it makes your code secure to share.

## Installation

```bash
go install github.com/sparkle-tech/obfuscator@latest
```

You can also build from source:

```bash
git clone https://github.com/devOpifex/obfuscator.git
cd obfuscator
go build
```

## Usage

```
$> obfuscator -h
Usage of obfuscator:
  -deobfuscate
        Deobfuscate the obfuscated files
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

### Examples

**Basic Obfuscation:**

```bash
obfuscator -in=R -out=obfuscated -key=123
```

**With License and Protected Tokens:**

```bash
obfuscator -in=R -out=obfuscated -key=123 -license=license.txt -protect=myFunction,importantVar
```

**Deobfuscation:**

```bash
obfuscator -deobfuscate -in=obfuscated -out=deobfuscated -key=123
```

### Parameter Details

- **-in**: Source directory containing R files to process
- **-out**: Destination directory for processed files
- **-key**: Encryption key used for the obfuscation algorithm
- **-license**: Path to a text file containing license information to add to each file
- **-protect**: Comma-separated list of identifiers that should not be obfuscated
- **-deobfuscate**: Flag to reverse the obfuscation process

## Limitations and Caveats

### Code Structure Requirements

- **If statements** must include curly braces:
  - ❌ `if(TRUE) 1 else 0`
  - ✅ `if(TRUE) {1} else {0}`

- **Lambda functions** must include curly braces:
  - ❌ `\(x) x + 1`
  - ✅ `\(x) {x + 1}`

- **Expressions in curly braces** outside of standard control structures are not supported:
  - ❌ `tryCatch({x + 1})`
  - ✅ `fn <- \(){x + 1}; tryCatch(fn())`

### Obfuscation Exceptions

- Functions starting with a dot (e.g., `.onLoad`) are **not obfuscated**
- File names starting with `__` are **not renamed** (but their content is still obfuscated)
- Only files witht the `.R` extension are processed
- Arguments to `do.call()` are **not obfuscated** - consider alternatives

### Best Practices

- Use a consistent and secure key for obfuscation/deobfuscation
- Keep a backup of your original code
- Test obfuscated code thoroughly before distribution
- Use the `-protect` flag for functions that must maintain their original names

## How It Works

The obfuscator works by:
1. Parsing R code to identify variables, functions, and other identifiers
2. Generating obfuscated names based on the provided key
3. Consistently replacing identifiers throughout the codebase
4. Preserving the functionality while making the code harder to read

