library(ambiorix)

source("./test2.R")

foo(41)

app <- Ambiorix$new(port = 8000L)

app$get("/", \(req, res) {
  res$send("Hello, World!")
})

app$start()
