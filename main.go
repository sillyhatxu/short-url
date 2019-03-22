package main

import (
	"flag"
	"short-url/bolt"
	"short-url/conf"
	"short-url/web"
)

func init() {
	dbName := flag.String("d", "short-url.db", "configuration file")
	schema := flag.String("s", "http", "configuration file")
	domainName := flag.String("dn", "127.0.0.1:8080", "configuration file")
	flag.Parse()
	conf.InitialConfig("data/"+*dbName, *schema, *domainName)
}

func main() {
	boltClient := bolt.NewBoltClient(conf.Conf.DBName, 0600)
	err := boltClient.InitialBucket()
	if err != nil {
		panic(err)
	}
	web.Start()
}
