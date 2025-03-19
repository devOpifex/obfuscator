foo <- \(x) {
  x <- x + 1L
  return(x)
}

foo(42)

x <- 2L

x <<- x + 1L

if (x > 1L) TRUE else FALSE

for (i in 1:10) {
  x <- x + 1L
}
