package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DimasPradana/kantor/payPDL-api/config"
	"github.com/DimasPradana/kantor/payPDL-api/database"
	modelpbb "github.com/DimasPradana/kantor/payPDL-api/model/pbb"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

var (
	resPelayanan modelpbb.StructResPelayanan
	arrPelayanan []modelpbb.StructResPelayanan
)

func Pelayanan(c *gin.Context) {
	/**
	TODO snub on Sen 12 Jul 2021 20:00:24 :
	- parsing ✓
	- query ✗
	- denda, total ✗
	*/

	_nop := c.Param("nop")
	kdKecamatan := _nop[4:7]
	kdKelurahan := _nop[7:10]
	kdBlok := _nop[10:13]
	noUrut := _nop[13:17]
	kdJnsOp := _nop[17:18]

	qry := fmt.Sprintf("select a.NM_WP_SPPT as nama, a.THN_PAJAK_SPPT as tahun, "+
		"a.PBB_YG_HARUS_DIBAYAR_SPPT as pokok, a.JLN_WP_SPPT as alamatwp, "+
		"a.STATUS_PEMBAYARAN_SPPT as lunas, a.NJOP_BUMI_SPPT as njopbumi, "+
		"a.NJOP_BNG_SPPT as njopbng, a.TGL_JATUH_TEMPO_SPPT as jatuhtempo, b.TGL_PEMBAYARAN_SPPT as tanggalbayar, "+
		"c.TOTAL_LUAS_BUMI as luasinduk, d.KD_ZNT as kode "+
		"from pbb.SPPT a left join pbb.PEMBAYARAN_SPPT b "+
		"on a.KD_KECAMATAN = b.KD_KECAMATAN and a.KD_KELURAHAN = b.KD_KELURAHAN and a.KD_BLOK = b.KD_BLOK and "+
		"a.NO_URUT = b.NO_URUT and a.KD_JNS_OP = b.KD_JNS_OP and a.THN_PAJAK_SPPT = b.THN_PAJAK_SPPT "+
		"left join pbb.DAT_OBJEK_PAJAK c "+
		"on a.KD_KECAMATAN = c.KD_KECAMATAN and a.KD_KELURAHAN = c.KD_KELURAHAN and a.KD_BLOK = c.KD_BLOK and "+
		"a.NO_URUT = c.NO_URUT and a.KD_JNS_OP = c.KD_JNS_OP "+
		"left join pbb.DAT_OP_BUMI d "+
		"on a.KD_KECAMATAN = d.KD_KECAMATAN and a.KD_KELURAHAN = d.KD_KELURAHAN and a.KD_BLOK = d.KD_BLOK and "+
		"a.NO_URUT = d.NO_URUT and a.KD_JNS_OP = d.KD_JNS_OP "+
		"where a.KD_KECAMATAN = '%v' "+
		"and a.KD_KELURAHAN = '%v' "+
		"and a.KD_BLOK = '%v' "+
		"and a.NO_URUT = '%v' "+
		"and a.KD_JNS_OP = '%v' "+
		"order by a.THN_PAJAK_SPPT asc", kdKecamatan, kdKelurahan, kdBlok, noUrut, kdJnsOp)

	// qry := fmt.Sprintf("select a.NM_WP_SPPT as nama, a.THN_PAJAK_SPPT as tahun, "+
	// "a.PBB_YG_HARUS_DIBAYAR_SPPT as pokok, a.JLN_WP_SPPT as Alamatwp "+
	// "from PBB.SPPT a where a.KD_KECAMATAN='%v' and a.KD_KELURAHAN='%v' and a.KD_BLOK='%v' "+
	// "and a.NO_URUT='%v' and a.KD_JNS_OP='%v'", kdKecamatan, kdKelurahan,kdBlok,noUrut, kdJnsOp)

	/*
		TODO snub on Sab 04 Jul 2020 11:11:28  : ambil config dari env file
	*/
	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	/*
		TODO snub on Sab 04 Jul 2020 11:11:48  : konek pake env file
	*/
	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)
	logrus.Infof("row 74 inquiry request nop : %v", _nop)

	rows, err := kon.Query(qry)
	if err != nil {
		logrus.Infof("errornya di baris 79 : %v,", qry)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	for rows.Next() {
		if err != rows.Scan(&resPelayanan.NAMA, &resPelayanan.TAHUN, &resPelayanan.POKOK, &resPelayanan.ALAMATWP,
			&resPelayanan.LUNAS, &resPelayanan.NJOPBUMI, &resPelayanan.NJOPBNG, &resPelayanan.JATUHTEMPO,
			&resPelayanan.TANGGALBAYAR, &resPelayanan.LUASINDUK, &resPelayanan.KODE) {
			// if err != rows.Scan(&resPelayanan.NAMA, &resPelayanan.TAHUN) {
			// if err != rows.Scan(&resPelayanan.NAMA, &resPelayanan.TAHUN, &resPelayanan.POKOK, &resPelayanan.ALAMATWP) {
			//http.Error(c.Writer, err.Error(), 500)
			//return
		}

		t := resPelayanan.JATUHTEMPO[0:10]
		_t, err := time.Parse("2006-01-02", t)
		if err != nil {
			logrus.Error(err)
			return
		}
		resPelayanan.DENDA = uint64(ambilDenda(_t, resPelayanan.POKOK))
		resPelayanan.TOTAL = resPelayanan.DENDA + resPelayanan.POKOK

		arrPelayanan = append(arrPelayanan, resPelayanan)

	}
	defer rows.Close()
	defer kon.Close()

	hasil := arrPelayanan
	c.JSON(http.StatusOK, hasil)
	arrPelayanan = nil
}
