package main

import (
	"github.com/reivaj05/url_shortener/shortener"

	"github.com/reivaj05/GoConfig"
	"github.com/reivaj05/GoLogger"
	"github.com/reivaj05/GoServer"
)

const (
	processName = "url-shortener"
)

func main() {
	setup()
	process()
}

func setup() {
	startConfig()
	startLogger()
}

func startConfig() {
	if err := GoConfig.Init(createConfigOptions()); err != nil {
		finishExecution("Error while loading config", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func createConfigOptions() *GoConfig.ConfigOptions {
	return &GoConfig.ConfigOptions{
		ConfigType: "json",
		ConfigFile: "config",
		ConfigPath: ".",
	}
}

func startLogger() {
	if err := GoLogger.Init(createLoggerOptions()); err != nil {
		finishExecution("Error while loading logger", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func createLoggerOptions() *GoLogger.LoggerOptions {
	return &GoLogger.LoggerOptions{
		OutputFile: processName + "-log.json",
		Path:       "log/",
		LogLevel:   getLogLevel(),
	}
}

func getLogLevel() int {
	levels := map[string]int{"DEBUG": GoLogger.DEBUG, "INFO": GoLogger.INFO,
		"WARNING": GoLogger.WARNING, "ERROR": GoLogger.ERROR,
		"PANIC": GoLogger.PANIC, "FATAL": GoLogger.FATAL,
	}
	if level, ok := levels[GoConfig.GetConfigStringValue("logLevel")]; ok {
		return level
	}
	return GoLogger.INFO
}

func process() {
	port := GoConfig.GetConfigStringValue("port")
	if err := GoServer.Start(port, createEndpoints()); err != nil {
		finishExecution("Error while starting server", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func createEndpoints() (endpoints []*GoServer.Endpoint) {
	endpoints = append(endpoints, shortener.Endpoints...)
	return
}

func finishExecution(msg string, fields map[string]interface{}) {
	GoLogger.LogFatal(msg, fields)
}
