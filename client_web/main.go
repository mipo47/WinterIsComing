package main

func main()  {
	httpServer := new(HttpServer)
	httpServer.Start(8080)
}
