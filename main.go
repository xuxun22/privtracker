package main

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var DOMAIN="tracker.whereami.com.cn"
var PORT="1337"

func main() {
	go Cleanup()
	config := fiber.Config{
		ServerHeader: DOMAIN,
		ReadTimeout:  time.Second * 245,
		WriteTimeout: time.Second * 15,
	}
	app := fiber.New(config)
	app.Use(recover.New())
	app.Use(myLogger())
	app.Static("/", "docs", fiber.Static{
		MaxAge:        3600 * 24 * 7,
		Compress:      true,
		CacheDuration: time.Hour,
	})
	app.Get("/monitor", monitor.New())
	app.Get("/:room/announce", announce)
	app.Get("/:room/scrape", scrape)
	app.Server().LogAllErrors = true
	log.Fatal(app.Listener(myListener()))
}

func myListener() net.Listener {
	var certDir="cert"
	cert, err := tls.LoadX509KeyPair(certDir+"/"+DOMAIN+".pem", certDir+"/"+DOMAIN+".key")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":" + PORT, cfg)
	if err != nil {
		panic(err)
	}
	return ln
}

func myLogger() fiber.Handler {
	loggerConfig := logger.ConfigDefault
	loggerConfig.Format = "${status} - ${latency} ${ip} ${method} ${path} ${bytesSent} - ${referer} - ${ua}\n"
	return logger.New(loggerConfig)
}

