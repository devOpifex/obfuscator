# this is a comment
# this is a comment
# this is a comment
# this is a comment
# this is a comment
get_home <- \(req, res) {
  x <- list(title = "Hello, World!")
  x$title |> res$send()
}
