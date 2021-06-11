package controller

import (
	"encoding/json"
	"fmt"
	"github.com/DimasPradana/kantor/payPDL-api/config"
	"github.com/DimasPradana/kantor/payPDL-api/database"
	"github.com/DimasPradana/kantor/payPDL-api/model/pbb"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var tagihan model.InqTagihan
var arrTagihan []model.InqTagihan
var reqInquiry model.StructReqInquiry
var resInquiry model.StructResInquiry
var status model.InqStatus
var InqError model.InqError
var InqStatusError model.InqStatusError
var nama model.InqNama
var arrNama []model.InqNama
var denda float64
var lunas []byte
var jatuhtempo string

func Inquiry(c *gin.Context) {

	/**
	TODO snub on Jum 05 Jun 2021 00:30:24 :
	- ketika gonta ganti NOP, Nama tetap tapi tagihan berubah ✓
	- pembayaran denda tidak masuk sismiop → PEMBAYARAN_SPPT ✗
	*/

	/*
		TODO snub on Sab 04 Jul 2020 11:10:49  : baca json body
	*/
	err := json.NewDecoder(c.Request.Body).Decode(&reqInquiry)
	if err != nil {
		logrus.Infof("errornya di baris 44: %v", err.Error())
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	/*
		TODO snub on Sab 04 Jul 2020 11:11:01  : parsing NOP
	*/
	kdKecamatan := reqInquiry.NOP[4:7]
	kdKelurahan := reqInquiry.NOP[7:10]
	kdBlok := reqInquiry.NOP[10:13]
	noUrut := reqInquiry.NOP[13:17]
	kdJnsOp := reqInquiry.NOP[17:18]
	tahun := reqInquiry.MASAPAJAK
	nop := kdKecamatan + "-" + kdKelurahan + "-" + kdBlok + "-" + noUrut + "-" + kdJnsOp

	/*
		TODO snub on Sab 04 Jul 2020 11:11:16  : query ke database masuk variabel
	*/

	qryFix := fmt.Sprintf("select * from (select a.THN_PAJAK_SPPT as Tahun, "+
		"a.PBB_YG_HARUS_DIBAYAR_SPPT as Pokok, a.TGL_JATUH_TEMPO_SPPT as jatuhtempo, "+
		"a.NM_WP_SPPT as nama, a.STATUS_PEMBAYARAN_SPPT as lunas, "+
		"b.NM_KELURAHAN as kelurahan "+
		"from sppt a "+
		"left join REF_KELURAHAN b on a.KD_KECAMATAN = b.KD_KECAMATAN and a.KD_KELURAHAN = b.KD_KELURAHAN "+
		"where a.KD_KECAMATAN = '%v' "+
		"and a.KD_KELURAHAN = '%v' "+
		"and a.KD_BLOK = '%v' "+
		"and a.NO_URUT = '%v' "+
		"and a.KD_JNS_OP = '%v' "+
		"and a.thn_pajak_sppt <= '%v' "+
		"and a.STATUS_PEMBAYARAN_SPPT = 0 "+
		"group by a.THN_PAJAK_SPPT, a.PBB_YG_HARUS_DIBAYAR_SPPT, a.TGL_JATUH_TEMPO_SPPT, a.NM_WP_SPPT, "+
		"a.STATUS_PEMBAYARAN_SPPT, b.NM_KELURAHAN "+
		"order by a.THN_PAJAK_SPPT desc) where ROWNUM <= 11", kdKecamatan, kdKelurahan, kdBlok, noUrut, kdJnsOp, tahun)

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
	// logrus.Infof("kon : %v", kon)
	logrus.Infof("row 95 inquiry request nop : %v", nop)

	rows, err := kon.Query(qryFix)
	if err != nil {
		logrus.Infof("errornya di baris 99 : %v,", qryFix)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	for rows.Next() {
		//if err != rows.Scan(&tagihan.TAHUN, &tagihan.POKOK, &tagihan.JATUHTEMPO, &nama.NAMA, &tagihan.LUNAS, &nama.KELURAHAN) {
		if err != rows.Scan(&tagihan.TAHUN, &tagihan.POKOK, &jatuhtempo, &nama.NAMA, &lunas, &nama.KELURAHAN) {
			http.Error(c.Writer, err.Error(), 500)
			return
		}

		t := jatuhtempo[0:10]
		_t, err := time.Parse("2006-01-02", t)
		if err != nil {
			logrus.Error(err)
			return
		}

		tagihan.DENDA = uint64(ambilDenda(_t, tagihan.POKOK))
		tagihan.TOTAL = tagihan.POKOK + tagihan.DENDA

		arrTagihan = append(arrTagihan, tagihan)
		arrNama = append(arrNama, nama)
		// logrus.Infof("isi dari scan\nNama : %v,\nTahun: %v\nTagihan: %v", arrNama[0].NAMA, tagihan.TAHUN, tagihan.POKOK)
		// logrus.Infof("\narray nama : %v", arrNama)
	}
	//logrus.Infof("isi dari scan\nNama : %v,\nPokok: %v\nDenda: %v\nTotal: %v", arrNama[0].NAMA, arrTagihan[2].POKOK, arrTagihan[2].DENDA, arrTagihan[2].TOTAL)

	/*
		TODO snub on Sab 04 Jul 2020 11:13:59  : inisialisasi respon body
	*/
	switch {
	//case arrNama[0].NAMA == "":
	case arrNama == nil:
		{
			InqError.ISERROR = "True"
			InqError.RESPONSECODE = "10"
			InqError.ERRORDESC = "Data SPPT dengan Tahun pajak tersebut tidak terdapat pada database"
			InqStatusError.STATUSERROR = InqError
			arrTagihan, arrNama = nil, nil
			c.JSON(http.StatusOK, InqStatusError)
		}
	case lunas[0] == 1:
		{
			InqError.ISERROR = "True"
			InqError.RESPONSECODE = "13"
			InqError.ERRORDESC = "Tagihan SPPT dengan Tahun Pajak dimaksud telah dibayar"
			InqStatusError.STATUSERROR = InqError
			arrTagihan, arrNama = nil, nil
			c.JSON(http.StatusOK, InqStatusError)
		}
	default:
		{
			resInquiry.NOP = reqInquiry.NOP
			resInquiry.NAMA = arrNama[0].NAMA
			resInquiry.KELURAHAN = arrNama[0].KELURAHAN
			resInquiry.KODEKP = "belum ada isi"
			resInquiry.KODEINSTITUSI = reqInquiry.KODEINSTITUSI
			resInquiry.NOHP = reqInquiry.NOHP
			resInquiry.EMAIL = reqInquiry.EMAIL
			resInquiry.TAGIHAN = arrTagihan
			status.ISERROR = "False"
			status.RESPONSECODE = "00"
			status.ERRORDESC = "Success"
			resInquiry.STATUS = status
			arrTagihan, arrNama = nil, nil
			c.JSON(http.StatusOK, resInquiry)
		}
	}
	defer kon.Close()
}

func ambilDenda(jatuhtempo time.Time, pokok uint64) float64 {
	duration := time.Since(jatuhtempo)
	hasil := math.Round(float64(int(duration.Hours()/24)) / 30)

	switch {
	case hasil <= 0:
		denda = 0
	case hasil >= 1 && hasil <= 24:
		denda = math.Round(0.02 * hasil * float64(pokok))
	default:
		denda = math.Round(0.48 * float64(pokok))
	}

	return denda
}
