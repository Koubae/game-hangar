package settings

import "time"

const AuthTokenExpirationTime = time.Hour * 4

const EnvTesting = "testing"
const EnvDev = "development"
const EnvStaging = "staging"
const EnvProd = "production"

var Environments = [4]string{EnvTesting, EnvDev, EnvStaging, EnvProd}
