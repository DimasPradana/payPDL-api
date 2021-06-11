package controller

import (
	"encoding/json"
	"fmt"
	"github.com/DimasPradana/kantor/payPDL-api/config"
	"github.com/DimasPradana/kantor/payPDL-api/database"
	"github.com/DimasPradana/kantor/payPDL-api/model/pbb"
	"net/http"
	"os"
	// "strings"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Reversal(c *gin.Context) {

	/*
		TODO snub on Rab 19 Agu 2020 11:27:22  :
		- di spek teknis response code 34 nomor 2 dan code 36 bedanya?
		- ketika rev_flag jadi 1, pembayaran_sppt di schema PBB hilang
		- response code belum diisi, default masih sukses
		- koneksi pake 1 kali saja
	*/

	var reqReversal model.StructReqReversal
	var resReversal model.StructResReversal
	var status model.RevStatus
	var RevError model.RevError
	var RevStatusError model.RevStatusError
	var thnpjk []string
	var countPay byte

	err := json.NewDecoder(c.Request.Body).Decode(&reqReversal)
	if err != nil {
		logrus.Infof("errornya di baris 38: %v", err.Error())
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	kdKecamatan := reqReversal.NOP[4:7]
	kdKelurahan := reqReversal.NOP[7:10]
	kdBlok := reqReversal.NOP[10:13]
	noUrut := reqReversal.NOP[13:17]
	nop := kdKecamatan + "-" + kdKelurahan + "-" + kdBlok + "-" + noUrut
	logrus.Infof("reversal request body : %v", nop)

	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)

	qry := fmt.Sprintf("select count(a.KD_PROPINSI) as countPay "+
		"from PBB.PEMBAYARAN_SPPT a "+
		"where a.KD_KECAMATAN = '%v' "+
		"and a.KD_KELURAHAN = '%v' "+
		"and a.KD_BLOK = '%v' "+
		"and a.NO_URUT = '%v'", kdKecamatan, kdKelurahan, kdBlok, noUrut)

	rows, err := kon.Query(qry)
	if err != nil {
		logrus.Infof("errornya di baris 68: %v\n%v", err.Error(), qry)
		// kesalahan = err.Error()
		http.Error(c.Writer, err.Error(), 500)
		return
	} else {
		/*
			TODO snub on Sab 04 Jul 2020 11:13:28  : ambil data dari database lalu masukkan ke dalam variabel
		*/
		for rows.Next() {
			// if err := rows.Scan(&nama.NAMA, &nama.KELURAHAN, &tagihan.TAHUN, &tagihan.POKOK, &tagihan.DENDA, &tagihan.TOTAL); err != nil {
			if err := rows.Scan(&countPay); err != nil {
				logrus.Infof("errornya di baris 79: %v", err.Error())
				// kesalahan = err.Error()
				http.Error(c.Writer, err.Error(), 500)
				return
				//} else {
				//	logrus.Infof("countPay: %v", countPay)
			} //}
		}
	}

	for a := range reqReversal.TAGIHAN {
		thnpjk = append(thnpjk, reqReversal.TAGIHAN[a].TAHUN)
		// thnpjkjoin = strings.Join(thnpjk, ",")
	}

	// logrus.Infof("thnpjk: %s", thnpjk)

	switch {
	case countPay == 0:
		{
			// resReversal.NOP = reqReversal.NOP
			// resReversal.REFERENCE = reqReversal.REFERENCE
			// status.ISERROR = "True"
			// status.RESPONSECODE = "34"
			// status.ERRORDESC = "Reversal Data Not Found"
			// resReversal.STATUS = status
			// c.JSON(http.StatusOK, resReversal)

			RevError.ISERROR = "True"
			RevError.RESPONSECODE = "13"
			RevError.ERRORDESC = "Tagihan SPPT dengan Tahun Pajak dimaksud telah dibayar"
			RevStatusError.STATUSERROR = RevError
			// logrus.Infof("testis: %v", InqError)
			c.JSON(http.StatusOK, RevStatusError)
		}
	default:
		{
			// logrus.Infof("panjang tahun pajak: %v,\n isinya: %v", len(spptThnPjk), spptThnPjk)
			for i := range thnpjk {
				//yoman := UpdateReversal(kdKecamatan, kdKelurahan, kdBlok, noUrut, thnpjk[i])
				UpdateReversal(kdKecamatan, kdKelurahan, kdBlok, noUrut, thnpjk[i])
				//logrus.Infof("\nReversalSPO: %v", yoman)
				// InsertPayment(kdKecamatan, kdKelurahan, kdBlok, noUrut, spptThnPjk[i], reqPayment.DATETIME, strings.Join(formatKodepengesahan, ""), spptPokok[i], spptDenda[i])
			}
			resReversal.NOP = reqReversal.NOP
			resReversal.REFERENCE = reqReversal.REFERENCE
			status.ISERROR = "False"
			status.RESPONSECODE = "00"
			status.ERRORDESC = "Success"
			resReversal.STATUS = status
			c.JSON(http.StatusOK, resReversal)
		}
	}

	defer kon.Close()
}

func UpdateReversal(kec, kel, blok, urut, thn string) string {

	queryDelSPO := fmt.Sprintf("delete "+
		"from SPO.PEMBAYARAN_SPPT "+
		"where KD_KECAMATAN = '%v' "+
		"and KD_KELURAHAN= '%v' "+
		"and KD_BLOK = '%v' "+
		"and NO_URUT = '%v' "+
		"and THN_PAJAK_SPPT = '%v'", kec, kel, blok, urut, thn)
	querydelPay := fmt.Sprintf("delete "+
		"from PBB.PEMBAYARAN_SPPT "+
		"where KD_KECAMATAN = '%v' "+
		"and KD_KELURAHAN= '%v' "+
		"and KD_BLOK = '%v' "+
		"and NO_URUT = '%v' "+
		"and THN_PAJAK_SPPT = '%v'", kec, kel, blok, urut, thn)
	queryUpdateSPPT := fmt.Sprintf("UPDATE PBB.SPPT b "+
		"SET b.STATUS_PEMBAYARAN_SPPT = 0 "+
		"WHERE b.KD_PROPINSI = '35' "+
		"AND b.KD_DATI2 = '12' "+
		"AND b.KD_KECAMATAN = '%s' "+
		"AND b.KD_KELURAHAN = '%s' "+
		"AND b.KD_BLOK = '%s' "+
		"AND b.NO_URUT = '%s' "+
		"AND b.THN_PAJAK_SPPT = '%s'", kec, kel, blok, urut, thn)

	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)
	hasil, err := kon.Exec(queryDelSPO)
	hasil3, err := kon.Exec(querydelPay)
	hasil2, err := kon.Exec(queryUpdateSPPT)
	if err != nil {
		logrus.Fatalf("errornya di baris 174: %v\n%v\nhasilnya: %v", err.Error(), queryDelSPO, hasil)
		logrus.Fatalf("errornya di baris 175: %v\n%v\nhasilnya: %v", err.Error(), querydelPay, hasil3)
		logrus.Fatalf("errornya di baris 176: %v\n%v\nhasilnya: %v", err.Error(), queryUpdateSPPT, hasil2)
	} else {
		logrus.Infof("Reversal %v-%v-%v-%v|%v Sukses", kec, kel, blok, urut, thn)
	}

	return queryDelSPO
}
