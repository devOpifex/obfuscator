box::use(. / test2[poolClose], pool[poolClose])

foo(23)

.onload <- function() {
  print("onload")
}

generic <- function(x) {
  UseMethod("generic")
}

generic.default <- function(x) {
  print("default")
}

generic.character <- function(x) {
  print("character")
}

generic.numeric <- function(x) {
  print("numeric")
}

generic.data.frame <- function(x) {
  print("data.frame")
}
