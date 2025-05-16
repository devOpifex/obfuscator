foo <- \(x) {
  x + 1
}

bar <- \(x) {
  foo(x)
}

baz <- \(x) {
  bar(x)
}

baz(42)


