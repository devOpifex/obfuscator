#' @export
foo <- \(x) {
  return(x + 1)
}

#' @export
`%||%` <- \(x, y) {
  if (is.null(x)) {
    y
  } else {
    x
  }
}
