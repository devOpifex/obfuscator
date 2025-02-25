#' @export
parse_str_to_date <- \(x, default) {
  is_empty <- is.null(x) || identical(length(x), 0L)
  if (is_empty) {
    return(default)
  }

  expr <- \(){
    res <- as.Date(x)
    is_empty <- identical(length(res), 0L)
    if (is_empty) {
      return(default)
    }

    res
  }

  tryCatch(
    expr = expr(),
    error = \(e) {
      default
    }
  )
}
