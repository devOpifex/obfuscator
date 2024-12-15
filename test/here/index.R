# this is a comment
# this is a comment
# this is a comment
# this is a comment
# this is a comment
#' @export
get_home <- \(req, res) {
  x <- list(title = "Hello, World!")
  x$title |> res$send()
}

#' @export
p_rint <- \(x){
  print(x)
}
