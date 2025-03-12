#' @export
box::use(
  . /
    env_vars[
      get_env_var,
      in_prod_mode,
      in_debug_mode,
    ],
  . /
    datetime[
      current_date,
      format_datetime,
      current_datetime,
      parse_str_to_date,
    ],
  . / rename[rename],
  . /
    operators[
      `%||%`,
      get_value_or_empty,
    ],
  . / parse_req[parse_req],
)

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
