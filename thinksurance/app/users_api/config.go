package users_api

type Config struct {
	Port int

	ServerKeepAliveEnforcementMinTime    int
	ServerKeepAlivePermitWithoutStream   bool
	ServerKeepAliveMaxConnectionIdle     int
	ServerKeepAliveMaxConnectionAge      int
	ServerKeepAliveMaxConnectionAgeGrace int
	ServerKeepAliveTime                  int
	ServerKeepAliveTimeout               int

	RepositoryFileDirectory string
}
