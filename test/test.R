box::use(
  ambiorix[Ambiorix],
  ./here[get_home, p_rint],
)

p_rint("starting...")

app <- Ambiorix$new(port = 8000L)

app$get("/", get_home)

app$start()
