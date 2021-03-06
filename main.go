package main

import (
	"runtime"

	"github.com/DimasPradana/kantor/payPDL-api/controller"
	controllerbphtb "github.com/DimasPradana/kantor/payPDL-api/controller/bphtb"
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
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
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

	rootRoutes := router.Group("/")
	{
		rootRoutes.GET("testing", controller.Testing)
		rootRoutes.POST("inquiry", controllerpbb.Inquiry)
		rootRoutes.POST("payment", controllerpbb.Payment)
		rootRoutes.POST("reversal", controllerpbb.Reversal)
	}
	bphtbRoutes := router.Group("/bphtb")
	{
		bphtbRoutes.GET("/cekprogresif/:noidentitasth", controllerbphtb.CekProgresif)
	}
	pbbRoutes := router.Group("/pbb")
	{
		pbbRoutes.GET("/pelayanan/:nop", controllerpbb.Pelayanan)
	}
	router.Run(":8002")
}
