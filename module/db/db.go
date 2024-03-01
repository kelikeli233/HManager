package db

type ConnectionArgs struct {
	MasterDSN              string
	ReplicaDSN             string
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int
}

func Connect(backend string, args ConnectionArgs) (*MySQLDB, error) {
	return connectMySQL(args)
}
