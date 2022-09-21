package drivers

type DBDump interface {
	Dump(dumpFile string) error
}

type DBDumper struct {
	DBName   string
	Username string
	Password string
	Host     string
	Port     int
}

func NewDBDumper() *DBDumper {
	return &DBDumper{
		DBName:   "temp",
		Username: "root",
		Password: "secret",
		Host:     "localhost",
		Port:     3306,
	}
}
