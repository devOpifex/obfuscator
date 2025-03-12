library(shiny)

footer <- function(x) {
  return(x)
}

tags$footer(footer(1))
tags$footer

z <- 23

x <- list(y = list(z = 1))

switch(x$y$z, "a" = 1, "b" = 2, "c" = 3)
