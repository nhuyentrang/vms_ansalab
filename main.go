package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	// db, err := Connect(
	// 	"localhost", // sqlHost
	// 	"5432",      // sqlPort
	// 	"vms",       // sqlDbName
	// 	"disable",   // sqlSslmode
	// 	"ansa",      // sqlUser
	// 	"ansa",      // sqlPassword
	// 	"public",    // currentSchema
	// )
	// if err != nil {
	// 	log.Fatalf("Error connecting to database: %v", err)
	// }
	// defer db.Close()

	// // Kiểm tra kết nối và thực hiện một truy vấn đơn giản
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatalf("Failed to ping the database: %v", err)
	// } else {
	// 	fmt.Println("Successfully connected to PostgreSQL!")
	// }

	// if !Connected {
	// 	log.Infof("SQL DB not configured, skipping")
	// } else {
	// 	err = Ping()
	// 	if err != nil {
	// 		log.Errorf(err.Error())
	// 	} else {
	// 		log.Info("Successfully connected to database")
	// 		// Migrate tables
	// 		err = Migrate(
	// 			&Event{},
	// 		)
	// 		if err != nil {
	// 			panic("Failed to AutoMigrate table! err: " + err.Error())
	// 		}
	// 	}
	// }
	// defer Close()

	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server CORE start")

	go HTTPAPIServer()
	go RTSPServer()
	go Storage.StreamChannelRunAll()

	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChanel
		log.WithFields(logrus.Fields{
			"module": "main",
			"func":   "main",
		}).Info("Server receive signal", sig)
		done <- true
	}()

	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server start success and waiting for signals")
	<-done
	Storage.StopAll()
	time.Sleep(2 * time.Second)
	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server stop working by signal")
}
