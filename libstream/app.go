package libstream

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	AppName = "stream"
)

var (
	DB     *sql.DB
	Logger = zap.S()
)

type Application struct {
	cfg    *viper.Viper
	Server *StreamServer

	configFile string

	listenAddr  string
	storageAddr string

	serverAPIEndpoint string

	storageName string
	storageUser string
	storagePass string
	timerValue  int

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
	app.rootCmd.PersistentFlags().StringVarP(&app.listenAddr, "service_address", "l", "localhost:8099", "service address")
	app.rootCmd.PersistentFlags().StringVarP(&app.storageAddr, "storage_address", "a", "localhost:8099", "DB address")
	app.rootCmd.PersistentFlags().StringVarP(&app.serverAPIEndpoint, "api", "p", "", "API URL endpoint")
	app.rootCmd.PersistentFlags().StringVarP(&app.storageName, "storage_name", "n", "", "DB name")
	app.rootCmd.PersistentFlags().StringVarP(&app.storageUser, "storage_user", "u", "", "DB user")
	app.rootCmd.PersistentFlags().StringVarP(&app.storagePass, "storage_password", "d", "", "DB password")
	app.rootCmd.PersistentFlags().IntVarP(&app.timerValue, "timer_value", "t", 100, "time to wait")
}

func (app *Application) InitConfig(configName, envPrefix string) {
	cfg := viper.New()

	cfg.SetEnvPrefix(envPrefix)
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg.SetDefault("server.addr", "localhost:8099")
	cfg.BindPFlag("server.addr", app.rootCmd.PersistentFlags().Lookup("service_address"))

	cfg.SetDefault("server.apiPrefix", "")
	cfg.BindPFlag("server.apiPrefix", app.rootCmd.PersistentFlags().Lookup("api"))

	cfg.SetDefault("storage.address", "")
	cfg.BindPFlag("storage.address", app.rootCmd.PersistentFlags().Lookup("storage_address"))

	cfg.SetDefault("storage.name", "stream")
	cfg.BindPFlag("storage.name", app.rootCmd.PersistentFlags().Lookup("storage_name"))

	cfg.SetDefault("storage.user", "testu")
	cfg.BindPFlag("storage.user", app.rootCmd.PersistentFlags().Lookup("storage_user"))

	cfg.SetDefault("storage.password", "testup")
	cfg.BindPFlag("storage.password", app.rootCmd.PersistentFlags().Lookup("storage_password"))

	cfg.SetDefault("root.token", "!csdf!25")

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

	u := app.cfg.GetString("storage.user")
	p := app.cfg.GetString("storage.password")
	n := app.cfg.GetString("storage.name")

	if db, ok := app.ConnectToDB(u, p, n); ok {
		DB = db
	}

	app.listenAddr = app.cfg.GetString("server.addr")
	app.Server = NewServer(ServerConfig{
		address:   app.cfg.GetString("server.addr"),
		rootToken: app.cfg.GetString("root.token"),
		apiPrefix: app.cfg.GetString("server.apiPrefix"),
	})

	app.Server.Timer = time.NewTimer(time.Second * time.Duration(app.timerValue))
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

	u := app.cfg.GetString("storage.user")
	p := app.cfg.GetString("storage.password")
	n := app.cfg.GetString("storage.name")

	if db, ok := app.ConnectToDB(u, p, n); ok {
		DB = db
	}

	//dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	//	app.cfg.GetString("storage.user"),
	//	app.cfg.GetString("storage.password"),
	//	app.cfg.GetString("storage.name"))
	//DB, err = sql.Open("postgres", dbinfo)
	//if err != nil {
	//	panic(err)
	//}

	app.listenAddr = app.cfg.GetString("server.addr")
	app.Server = NewServer(ServerConfig{
		address:   app.cfg.GetString("server.addr"),
		rootToken: app.cfg.GetString("root.token"),
		apiPrefix: app.cfg.GetString("server.apiPrefix"),
	})

	app.Server.Timer = time.NewTimer(time.Second * time.Duration(app.timerValue))
}

func (app *Application) ConnectToDB(user, password, nameDB string) (*sql.DB, bool) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=5432 sslmode=disable",
		user,
		password,
		nameDB)

	d, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, false
	}
	return d, true
}

func (app *Application) Run() {
	if err := app.rootCmd.Execute(); err != nil {
		panic(err)
	}
}
