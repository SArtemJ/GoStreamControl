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