package conf

type Config struct {
	DBName     string // short-url
	Schema     string // http
	DomainName string // www.demo.com
}

var Conf Config

func InitialConfig(DBName, Schema, DomainName string) {
	Conf = Config{DBName: DBName, Schema: Schema, DomainName: DomainName}
}
