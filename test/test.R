foo <- \(x) {
  x <- x + 1L
  return(x)
}

foo(42)

x <- 2L

x <<- x + 1L
