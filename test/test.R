library(ambiorix)

foo <- function(x) {
  return(x + 1)
}

foo(41)

app <- Ambiorix$new(port = 8000L)

app$get("/", \(req, res) {
  res$send("Hello, World!")
})

app$start()
