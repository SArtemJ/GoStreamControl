package libstream

import (
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	"github.com/gorilla/mux"
	"net/http"
)

//const (
//	APP_NAME string = "stream"
//)

type Application struct {

	cfg *viper.Viper

	interruptTimer int
	rootToken string

	rootCmd *cobra.Command
}

func NewApplication() *Application {
	app := Application{}
	return &app
}

func (app *Application) InitCommands() {

	app.rootCmd = &cobra.Command{
		Use: "stream",
		Short: "Stream API",
		Long: "Stream control API",
		Run: func(cmd *cobra.Command, args []string) {
			app.Init()
		},
	}

	app.rootCmd.PersistentFlags().IntVarP(&app.interruptTimer, "timer", "t", 1000, "")
	app.rootCmd.PersistentFlags().StringVarP(&app.rootToken, "rootToken", "rt", "aVfg!&afP" ,"")
}


func (app *Application) Init() {
	app.interruptTimer = app.cfg.GetInt("timer")
	app.rootToken = app.cfg.GetString("rootToken")

	router := mux.NewRouter()
	sub := router.PathPrefix("/api/v1").Subrouter()
	sub.HandleFunc("/s", ShowAllStreams).Methods("GET")
	sub.HandleFunc("/run", StartNewStream).Methods("GET")
	sub.HandleFunc("/activate/{id}", ActivateStream).Methods("PATCH")
	sub.HandleFunc("/interrupt/{id}", InterruptStream).Methods("PATCH")
	sub.HandleFunc("/finish/{id}", FinishStream).Methods("PATCH")
	http.ListenAndServe(":8000", router)

}


func (app *Application) Run() {
	if err := app.rootCmd.Execute(); err != nil {
		panic(err)
	}
}