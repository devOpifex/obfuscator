count_num_of_errors <- \(pool, n = 1L) {
  cur <- as.Date(current_date()) - n + 1L
  prev <- as.Date(current_date()) - n * 2L + 1L

  case_cur <- sprintf(
    "SUM(CASE WHEN event_time >= '%s' AND type = 'error' THEN 1 ELSE 0 END) AS cur_events",
    cur
  )

  case_prev <- sprintf(
    "SUM(CASE WHEN event_time >= '%s' AND event_time < '%s' AND type = 'error' THEN 1 ELSE 0 END) AS prev_events",
    prev, cur
  )

  case <- paste(case_cur, ", ", case_prev)

  query <- paste("SELECT", case, "FROM events")

  conn <- checkout_conn(pool = pool)
  on.exit(poolReturn(conn))

  res <- dbSendQuery(conn = conn, statement = query)
  found <- dbFetch(res)
  dbClearResult(res)

  found <- lapply(found, as.integer)

  pct_change <- paste0(
    round(
      (found$cur_events - found$prev_events) / found$prev_events * 100,
      digits = 1L
    ),
    "%"
  )

  if (identical(found$prev_events, 0L)) {
    pct_change <- "-"
  }

  status <- if (found$cur_events >= found$prev_events) {"increase"} else {"decrease"}

  found$pct_change <- pct_change
  found$status <- status

  # error rate
  all_cur_events <- count_events(pool = pool, n = n)[["cur_events"]]
  error_rate <- paste0(
    round(
      x = found$cur_events / all_cur_events * 100,
      digits = 1L
    ),
    "%"
  )

  if (identical(all_cur_events, 0L)) {
    error_rate <- "-"
  }

  found$error_rate <- error_rate

  found
}
