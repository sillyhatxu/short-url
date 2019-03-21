package main

import (
	"flag"
	"short-url/conf"
	"short-url/web"
)

func init() {
	dbName := flag.String("d", "short-url.db", "configuration file")
	schema := flag.String("s", "http", "configuration file")
	domainName := flag.String("dn", "127.0.0.1:8080", "configuration file")
	flag.Parse()
	conf.InitialConfig(*dbName, *schema, *domainName)
}

func main() {
	web.Start()
}
