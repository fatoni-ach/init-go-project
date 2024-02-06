package config

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	defCfg      map[string]string
	initialized = false
)

func initialize() {
	viper.SetConfigFile(`./config.yml`)

	viper.ReadInConfig()

	defCfg = make(map[string]string)

	defCfg["debug"] = "false"
	defCfg["key"] = "false"
	defCfg["name"] = "asset-service.com"

	defCfg["server.address"] = "8080"

	defCfg["context.timeout"] = "2"

	defCfg["database.host"] = "localhost"
	defCfg["database.port"] = "3306"
	defCfg["database.user"] = "devuser"
	defCfg["database.password"] = "devpassword"
	defCfg["database.name"] = "devdb"

	// defCfg["redis.host"] = "localhost"
	// defCfg["redis.port"] = "6379"

	// defCfg["email.username"] = ""
	// defCfg["email.password"] = ""
	// defCfg["email.address"] = "no-reply@idntimes.com"
	// defCfg["email.name"] = "Go Project"
	// defCfg["email.driver"] = "smtp"
	// defCfg["email.host"] = "10.100.100.107"
	// defCfg["email.port"] = "25"

	defCfg["cdn.url"] = "localhost"
	defCfg["cdn.oss"] = "false"

	defCfg["log.level"] = "trace"
	defCfg["log.path"] = "storage"
	defCfg["log.type"] = "file"
	defCfg["log.max.age"] = "10"

	defCfg["edit.count"] = "1"

	for k := range defCfg {
		err := viper.BindEnv(k)
		if err != nil {
			log.Errorf("Failed to bind env \"%s\" into configuration. Got %s", k, err)
		}
	}

	initialized = true
}

// SetConfig put configuration key value
func SetConfig(key, value string) {
	viper.Set(key, value)
}

// Get fetch configuration as string value
func Get(key string) string {
	if !initialized {
		initialize()
	}
	ret := viper.GetString(key)
	if len(ret) == 0 {
		if ret, ok := defCfg[key]; ok {
			return ret
		}
		log.Debugf("%s config key not found", key)
	}
	return ret
}

// GetBoolean fetch configuration as boolean value
func GetBoolean(key string) bool {
	if len(Get(key)) == 0 {
		return false
	}
	b, err := strconv.ParseBool(Get(key))
	if err != nil {
		panic(err)
	}
	return b
}

// GetInt fetch configuration as integer value
func GetInt(key string) int {
	if len(Get(key)) == 0 {
		return 0
	}
	i, err := strconv.ParseInt(Get(key), 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}

// GetFloat fetch configuration as float value
func GetFloat(key string) float64 {
	if len(Get(key)) == 0 {
		return 0
	}
	f, err := strconv.ParseFloat(Get(key), 64)
	if err != nil {
		panic(err)
	}
	return f
}

// Set configuration key value
func Set(key, value string) {
	defCfg[key] = value
}
