# this is a comment
# this is a comment
# this is a comment
# this is a comment
# this is a comment
get_home <- \(req, res) {
  x <- list(title = "Hello, World!")
  res$send(x$title)
}
