package controller

import (
	"fmt"
	"github.com/DimasPradana/kantor/payPDL-api/config"
	"github.com/DimasPradana/kantor/payPDL-api/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func CekProgresif(c *gin.Context) {
	var countProgresif *byte
	var namaPembeli, identitasTH *string
	noidentitasth := c.Param("noidentitasth")

	qry := fmt.Sprintf("select count(a.ID_PENGAJUAN), a.NAMA_TH, a.NO_IDENTITAS_TH "+
		" from bphtb_apps_bphtb.pengajuan a "+
		" where a.NO_IDENTITAS_TH like '%v%%' "+
		" group by a.NAMA_TH, a.NO_IDENTITAS_TH", noidentitasth)

	config.Getenvfile()
	envUser := os.Getenv("usermysql")
	envPass := os.Getenv("passwordmysql")
	envAddr := os.Getenv("addrmysql")
	envPort := os.Getenv("portmysql")

	kon, _ := database.KonekMysql(envAddr, envUser, envPass, envPort, "")
	rows, err := kon.Query(qry)
	if err != nil {
		//logrus.Infof("errornya di baris 32 : %v,", qry)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	for rows.Next() {
		if err != rows.Scan(&countProgresif, &namaPembeli, &identitasTH) {
			http.Error(c.Writer, err.Error(), 500)
			return
		}
	}
	defer rows.Close()
	defer kon.Close()

	hasil := fmt.Sprintf("Nomor KTP '%v' dengan nama '%v' telah melakukan pembelian sebanyak : %v kali dalam tahun ini", *identitasTH, *namaPembeli, *countProgresif)
	c.JSON(http.StatusOK, hasil)
    // identitasTH = ""; namaPembeli = ""; countProgresif = 0
}
