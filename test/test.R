box::use(
  ambiorix[Ambiorix],
  . / test2[foo, `%||%`],
)

NULL %||% foo(41)

app <- Ambiorix$new(port = 8000L)

app$get("/", \(req, res) {
  res$send("Hello, World!")
})

app$start()
