# LICENSE SPARKLE TECH
box::use(ambiorix[Ambiorix],./here[get_home,p_rint],);app=Ambiorix$new(port=8000L);app$get("/",get_home);app$start();