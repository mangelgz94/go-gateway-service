package find_number_position_api

type Config struct {
	Port int

	ServerKeepAliveEnforcementMinTime    int
	ServerKeepAlivePermitWithoutStream   bool
	ServerKeepAliveMaxConnectionIdle     int
	ServerKeepAliveMaxConnectionAge      int
	ServerKeepAliveMaxConnectionAgeGrace int
	ServerKeepAliveTime                  int
	ServerKeepAliveTimeout               int

	ArraySize int
}
