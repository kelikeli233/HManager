package db

import "ehmanager/module/datatypes"

type ConnectionArgs struct {
	Master                 datatypes.DBNodeConfig
	Replica                datatypes.DBNodeConfig
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int
}

func Connect(backend string, args ConnectionArgs) (*MySQLDB, error) {
	return connectMySQL(args)
}
