package libstream

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	AppName = "StreamApp"
)

var (
	Logger = zap.S()
)

type Application struct {
	cfg    *viper.Viper
	Server *StreamServer

	configFile        string
	listenAddr        string
	serverAPIEndpoint string
	storageAddr       string
	storageName       string
	timerValue        int

	Rt      string
	rootCmd *cobra.Command
}

func NewApplication() *Application {
	app := Application{}
	return &app
}

func (app *Application) InitCommands() {

	app.rootCmd = &cobra.Command{
		Use:   "stream",
		Short: "Stream API",
		Long:  "Stream control API",
		Run: func(cmd *cobra.Command, args []string) {
			app.Init()
			app.Server.Run()
		},
	}

	app.rootCmd.PersistentFlags().StringVarP(&app.configFile, "config", "c", "", "default ./libstream.yaml")
	app.rootCmd.PersistentFlags().StringVarP(&app.listenAddr, "service_address", "l", "localhost:8888", "Service address")
	app.rootCmd.PersistentFlags().StringVarP(&app.storageAddr, "storage_address", "a", "localhost", "Mongo address")
	app.rootCmd.PersistentFlags().StringVarP(&app.serverAPIEndpoint, "api", "p", "", "API URL endpoint")
	app.rootCmd.PersistentFlags().StringVarP(&app.storageName, "storage_name", "n", "Stream", "Mongo DB name")
	app.rootCmd.PersistentFlags().IntVarP(&app.timerValue, "timer_value", "t", 1, "Time to wait")
	app.rootCmd.PersistentFlags().StringVarP(&app.Rt, "root_token", "r", "", "Root token for delete")
}

func (app *Application) InitConfig(configName, envPrefix string) {
	cfg := viper.New()

	cfg.SetEnvPrefix(envPrefix)
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//on local-machine localhost:8888
	//on docker 0.0.0.0:8888
	cfg.SetDefault("server.addr", "0.0.0.0:8888")
	cfg.BindPFlag("server.addr", app.rootCmd.PersistentFlags().Lookup("service_address"))

	cfg.SetDefault("server.apiPrefix", "")
	cfg.BindPFlag("server.apiPrefix", app.rootCmd.PersistentFlags().Lookup("api"))

	//on local-machine localhost
	//on docker gostreamcontrolapi_db_1
	cfg.SetDefault("storage.address", "gostreamcontrolapi_db_1")
	cfg.BindPFlag("storage.address", app.rootCmd.PersistentFlags().Lookup("storage_address"))

	cfg.SetDefault("storage.name", "Stream")
	cfg.BindPFlag("storage.name", app.rootCmd.PersistentFlags().Lookup("storage_name"))

	cfg.SetDefault("timer.value", 1)
	cfg.BindPFlag("timer.value", app.rootCmd.PersistentFlags().Lookup("timer_value"))

	cfg.SetDefault("r.t", "!csdf!25")
	cfg.BindPFlag("r.t", app.rootCmd.PersistentFlags().Lookup("root_token"))

	cfg.BindPFlag("config", app.rootCmd.PersistentFlags().Lookup("config"))
	if cfg.GetString("config") != "" {
		cfg.SetConfigName(cfg.GetString("config"))
	} else {
		cfg.SetConfigName(configName)
	}

	cfg.AddConfigPath("/etc/")
	cfg.AddConfigPath("$HOME/")
	cfg.AddConfigPath("./")

	app.cfg = cfg
}

func (app *Application) GetConfig() *viper.Viper {
	return app.cfg
}

func (app *Application) ConfigureLog() {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger")
	}
	Logger = logger.Sugar()
}

func (app *Application) Configure(params ...string) {
	configName := AppName
	envName := AppName
	switch {
	case len(params) == 1:
		configName = params[0]
		envName = params[0]
	case len(params) > 1:
		configName = params[0]
		envName = params[1]
	}
	app.InitCommands()
	app.InitConfig(configName, envName)
	app.ConfigureLog()
}

func (app *Application) Init() {
	if app.configFile != "" {
		app.cfg.SetConfigName(app.configFile)
	}
	err := app.cfg.ReadInConfig()
	if err != nil {
		Logger.Debug("Configuration file not found")
	} else {
		Logger.Infow("Configuration file", "path", app.cfg.ConfigFileUsed())
	}

	app.listenAddr = app.cfg.GetString("server.addr")
	storage := NewMongoStorage(app.cfg.GetString("storage.address"), app.cfg.GetString("storage.name"))
	storage.Reset()

	app.Server = NewServer(ServerConfig{
		Address:    app.cfg.GetString("server.addr"),
		RootToken:  app.cfg.GetString("r.t"),
		ApiPrefix:  app.cfg.GetString("server.apiPrefix"),
		TimerValue: app.cfg.GetInt("timer.value"),
		Storage:    storage,
	})
}

func (app *Application) InitWithConfig(cfg map[string]interface{}) {
	if app.configFile != "" {
		app.cfg.SetConfigName(app.configFile)
	}
	err := app.cfg.ReadInConfig()
	if err != nil {
		Logger.Debug("Configuration file not found")
	} else {
		Logger.Infow("Configuration file", "path", app.cfg.ConfigFileUsed())
	}
	for key, value := range cfg {
		app.cfg.Set(key, value)
	}

	app.listenAddr = app.cfg.GetString("server.addr")
	storage := NewMongoStorage(app.cfg.GetString("storage.address"), app.cfg.GetString("storage.name"))
	storage.Reset()

	app.Server = NewServer(ServerConfig{
		Address:    app.cfg.GetString("server.addr"),
		RootToken:  app.cfg.GetString("r.t"),
		ApiPrefix:  app.cfg.GetString("server.apiPrefix"),
		TimerValue: app.cfg.GetInt("timer.value"),
		Storage:    storage,
	})
}

func (app *Application) Run() {
	if err := app.rootCmd.Execute(); err != nil {
		panic(err)
	}
}
