package libstream


var testApp *Application

func GetTestApp(cfg map[string]interface{}) *Application {
	if testApp == nil {
		testApp = NewApplication()
		testApp.Configure("stream_test")
		testApp.InitWithConfig(cfg)
	}
	return testApp
}

func GetTestServer() *StreamServer {
	return GetTestApp(nil).Server
}

func GetTestingServerWithConfig(cfg map[string]interface{}) *StreamServer {
	app := NewApplication()
	app.Configure("stream_test", "stream_test")
	app.InitWithConfig(cfg)
	return app.Server
}