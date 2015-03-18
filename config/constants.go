package config

// Current version
const CurrentVersion = "1"

// SessionTTLMinutes is a default TTL for all user sessions
const SessionTTLMinutes = 30

// Used for password hashing
//TODO: move to external configuration
const HashSecretSalt string = "5c8f28d559f89414e8D317f28850A32c"

// Context key to store router's params
const CtxParamsKey = "params"

// Context key to store session info
const CtxSessionKey = "session"
