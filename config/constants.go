package config

const (
	// CurrentAPIVersion is current IdP API version
	CurrentAPIVersion = "1"
	// CurrentCLIVersion is current IdP's CLI version
	CurrentCLIVersion = "0.0.1"

	// EnvIDPAddr environment variable
	EnvIDPAddr = "IDP_ADDR"
	// EnvIDPDriver environment variable
	EnvIDPDriver = "IDP_DB_Driver"
	// EnvIDPDSN environment variable
	EnvIDPDSN = "IDP_DB_DSN"
	// EnvIDPSessionTTL environment variable
	EnvIDPSessionTTL = "IDP_SESSION_TTL"
	// EnvIDPSecretSalt environment variable
	EnvIDPSecretSalt = "IDP_SECRET_SALT"
	// EnvIDPSQLTrace environment variable
	EnvIDPSQLTrace = "IDP_SQL_TRACE"

	// CtxParamsKey key to store router's params
	CtxParamsKey = "params"
	// CtxSessionKey key to store session info
	CtxSessionKey = "session"
)
