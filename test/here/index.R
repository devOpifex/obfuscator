# this is a comment
# this is a comment
# this is a comment
# this is a comment
# this is a comment
`%||%` <- function(x, y) {
  if (is.null(x)) y else x
}

get_home <- \(req, res) {
  1 %||% 2
  x <- list(title = "Hello, World!")
  res$send(x$title)
}
