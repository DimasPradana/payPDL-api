package main

import (
	"runtime"

	"github.com/DimasPradana/kantor/payPDL-api/controller"
	"github.com/gin-gonic/gin"
	// "io"
	controllerpbb "github.com/DimasPradana/kantor/payPDL-api/controller/pbb"
	// "os"
)

func main() {
	/*
		TODO snub on Min 23 Agu 2020 08:15:10  :
			- log file belum pake logrus
	*/
	runtime.GOMAXPROCS(4) //jumlah maksimal prosesor yang digunakan
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)
	// Force log's color
	// gin.ForceConsoleColor()
	// Disable Console Color, you don't need console color when writing the logs to file.
	// gin.DisableConsoleColor()

	// Logging to a file.
	// f, _ := os.Create("dimas.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	router.GET("/testing", controller.Testing)
	router.POST("/inquiry", controllerpbb.Inquiry)
	router.POST("/payment", controllerpbb.Payment)
	router.POST("/reversal", controllerpbb.Reversal)
	router.Run(":8002")
}
