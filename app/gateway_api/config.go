package gateway_api

type Config struct {
	Port int

	AuthUser     string
	AuthPassword string

	GRPCClientKeepAliveTime       int
	GRPCClientAliveTimeout        int
	GRPCClientPermitWithoutStream bool
	GRPCClientMAxAttempts         int
	GRPCClientMaxBackoff          string
	GRPCClientBackoffMultiplier   float64

	GRPCUsersAddress              string
	GRPCFindNumberPositionAddress string
}
