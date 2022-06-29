package app

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type RPCConfig struct {
	RPCOrderPort  string
	RPCUserPort   string
	RPCDriverPort string
	WaitingTime   int
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	SSLMode    string
}

func NewRPCConfig() (*RPCConfig, error) {
	res := &RPCConfig{}
	viper.AutomaticEnv()

	if !viper.IsSet("RPCORDERPORT") {
		return nil, errors.New("env var RPCORDERPORT is empty")
	}
	res.RPCOrderPort = viper.GetString("RPCORDERPORT")

	if !viper.IsSet("RPCUSERPORT") {
		return nil, errors.New("env var RPCUSERPORT is empty")
	}
	res.RPCUserPort = viper.GetString("RPCUSERPORT")

	if !viper.IsSet("RPCDRIVERPORT") {
		return nil, errors.New("env var RPCDRIVERPORT is empty")
	}
	res.RPCDriverPort = viper.GetString("RPCDRIVERPORT")

	if !viper.IsSet("WAITINGTIME") {
		return nil, errors.New("env var WAITINGTIME is empty")
	}
	res.WaitingTime = viper.GetInt("WAITINGTIME")

	return res, nil
}

func NewDBConfig() (*DBConfig, error) {
	res := &DBConfig{}
	viper.AutomaticEnv()

	if !viper.IsSet("DBUSER") {
		return nil, errors.New("env var DBUSER is empty")
	}
	res.DBUser = viper.GetString("DBUSER")

	if !viper.IsSet("DBPASSWORD") {
		return nil, errors.New("env var DBPASSWORD is empty")
	}
	res.DBPassword = viper.GetString("DBPASSWORD")

	if !viper.IsSet("DBHOST") {
		return nil, errors.New("env var DBHOST is empty")
	}
	res.DBHost = viper.GetString("DBHOST")

	if !viper.IsSet("DBPORT") {
		return nil, errors.New("env var DBPORT is empty")
	}
	res.DBPort = viper.GetString("DBPORT")

	if !viper.IsSet("DBNAME") {
		return nil, errors.New("env var DBNAME is empty")
	}
	res.DBName = viper.GetString("DBNAME")

	if !viper.IsSet("SSLMODE") {
		return nil, errors.New("env var SSLMODE is empty")
	}
	res.SSLMode = viper.GetString("SSLMODE")

	return res, nil
}
