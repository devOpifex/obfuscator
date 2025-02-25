#' Seed database with dummy data
#'
#' @details Adds dummy data to the DB if in debug mode.
#' Called during setup. See [. / conn[setup]].
#' @param pool The pool object.
#' @return `NULL`
seed_database <- \(pool) {
  events[,
    event_time := format_datetime(
      Sys.time() -
        sample(
          x = 1:(24 * 60 * 60 * 60),
          size = nrow(events),
          replace = TRUE
        )
    )
  ]
}
