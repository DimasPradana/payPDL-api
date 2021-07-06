package controller

import (
	"encoding/json"
	"fmt"
	"github.com/DimasPradana/kantor/payPDL-api/config"
	"github.com/DimasPradana/kantor/payPDL-api/database"
	"github.com/DimasPradana/kantor/payPDL-api/model/pbb"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Payment(c *gin.Context) {

	/*
		TODO snub on Sel 6 Juli Agu 2021 22:16:22  :
        - cek status lunas dan count di table bayar dulu sebelum payment ✗
        - jika status lunas dan count lebih dari 1 then beri respon lunas ✗
	*/

	var reqPayment model.StructReqPayment
	var resPayment model.StructResPayment
	var status model.PayStatus
	var PayError model.PayError
	var PayStatusError model.PayStatusError

	var formatTanggal = time.Now().Format("150405020106")
	var thnpjk []string
	var thnpjkjoin string
	var pokok, total, subtotal uint64
	var jatuhtempo string
	var denda float64
	// var arrPokok, arrDenda, arrTotal []uint64
	var arrDenda, arrTotal []uint64

	err := json.NewDecoder(c.Request.Body).Decode(&reqPayment)
	if err != nil {
		logrus.Infof("errornya di baris 45: %v", err.Error())
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	kdKecamatan := reqPayment.NOP[4:7]
	kdKelurahan := reqPayment.NOP[7:10]
	kdBlok := reqPayment.NOP[10:13]
	noUrut := reqPayment.NOP[13:17]
	kdJnsOp := reqPayment.NOP[17:18]
	nop := kdKecamatan + "-" + kdKelurahan + "-" + kdBlok + "-" + noUrut + "-" + kdJnsOp
	logrus.Infof("payment request body : %v", nop)

	/* kodepengesahan 16 digit tahun, tanggal,
	bulan, detik, jam, menit,
	urut1,urut2 */
	formatKodepengesahan := []string{formatTanggal[10:12], formatTanggal[6:8],
		formatTanggal[8:10], formatTanggal[4:6], formatTanggal[0:2], formatTanggal[2:4],
		noUrut[0:2], noUrut[2:4]}

	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)

	for a := range reqPayment.TAGIHAN {
		thnpjk = append(thnpjk, reqPayment.TAGIHAN[a].TAHUN)
		thnpjkjoin = strings.Join(thnpjk, ",")
		//logrus.Infof("reqPayment tagihan : %v", thnpjkjoin)
	}

	qryFix := fmt.Sprintf("select a.PBB_YG_HARUS_DIBAYAR_SPPT as pokok, a.TGL_JATUH_TEMPO_SPPT as jatuhtempo "+
		"from SPPT a "+
		"where a.KD_KECAMATAN = '%v' and a.KD_KELURAHAN = '%v' and a.KD_BLOK = '%v' and a.NO_URUT = '%v' and a.KD_JNS_OP = '%v' "+
		"and a.THN_PAJAK_SPPT in (%v)", kdKecamatan, kdKelurahan, kdBlok, noUrut, kdJnsOp, thnpjkjoin)

	rows, err := kon.Query(qryFix)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		logrus.Infof("errornya di baris 88 : %v", err.Error())
		return
	}

	for rows.Next() {
		if err := rows.Scan(&pokok, &jatuhtempo); err != nil {
			http.Error(c.Writer, err.Error(), 500)
			logrus.Infof("errornya di baris 95 : %v", err.Error())
			return
		}

		t := jatuhtempo[0:10]
		_t, err := time.Parse("2006-01-02", t)
		if err != nil {
			logrus.Error(err)
			return
		}
		denda = ambilDenda(_t, pokok)
		total = pokok + uint64(denda)

		// arrPokok = append(arrPokok, pokok)
		arrDenda = append(arrDenda, uint64(denda))
		arrTotal = append(arrTotal, total)
		//logrus.Infof("isi pokok : %v, jatuhtempo : %v, denda : %v, total : %v", pokok, jatuhtempo, denda, total)
	}
	//logrus.Infof("isi arrpokok : %v, arrdenda : %v, total : %v", arrPokok, arrDenda, arrTotal)
	defer rows.Close()

	for b := range arrTotal {
		subtotal += arrTotal[b]
		//logrus.Infof("total : %v", arrTotal[b])
	}
	//logrus.Infof("subtotal : %v", subtotal)

	switch {
	case arrTotal == nil:
		{
			//logrus.Infof("Tagihan Nihil")
			PayError.ISERROR = "True"
			PayError.RESPONSECODE = "23"
			PayError.ERRORDESC = "Jumlah Pembayaran Nihil"
			PayStatusError.STATUSERROR = PayError
			c.JSON(http.StatusOK, PayStatusError)
		}
	case reqPayment.TOTALBAYAR != subtotal:
		{
			logrus.Infof("Nominal tidak sesuai, harusnya: %v, yg diajukan: %v", subtotal, reqPayment.TOTALBAYAR)
			PayError.ISERROR = "True"
			PayError.RESPONSECODE = "16"
			PayError.ERRORDESC = "Jumlah Nominal yang dibayarkan tidak sama dengan seharusnya dibayar"
			PayStatusError.STATUSERROR = PayError
			c.JSON(http.StatusOK, PayStatusError)
		}
	default:
		{
			kodetp := ambilKodeTP(reqPayment.KODEINSTITUSI)
			//logrus.Infof("isi dari kodetp baris 143 : %v, kodebanktunggal : %v", kodetp, kodetp.KD_BANK_TUNGGAL)
			//for i := range thnpjk {
			//	InsertPayment(kdKecamatan, kdKelurahan, kdBlok, noUrut, kdJnsOp, thnpjk[i], kodetp.KD_KANWIL,
			//		kodetp.KD_KPPBB, kodetp.KD_BANK_TUNGGAL, kodetp.KD_BANK_PERSEPSI, kodetp.KD_TP,
			//		reqPayment.DATETIME, strings.Join(formatKodepengesahan, ""), reqPayment.KODEINSTITUSI, arrDenda[i], arrTotal[i])
			//}
			for i := range thnpjk {
				InsertPayment(kdKecamatan, kdKelurahan, kdBlok, noUrut, kdJnsOp, thnpjk[i], kodetp.KD_KANWIL,
					kodetp.KD_KPPBB, kodetp.KD_BANK_TUNGGAL, kodetp.KD_BANK_PERSEPSI, kodetp.KD_TP,
					reqPayment.DATETIME, strings.Join(formatKodepengesahan, ""), kodetp.KODE_INSTITUSI, arrDenda[i], arrTotal[i])
			}
			resPayment.NOP = reqPayment.NOP
			resPayment.KODEPENGESAHAN = strings.Join(formatKodepengesahan, "")
			resPayment.KODEKP = "0000"
			status.ISERROR = "False"
			status.RESPONSECODE = "00"
			status.ERRORDESC = "Success"
			resPayment.STATUS = status
			c.JSON(http.StatusOK, resPayment)
		}
	}

	defer kon.Close()
}

func InsertPayment(kec, kel, blok, urut, kdjnsop, thn, kdtp1, kdtp2, kdtp3, kdtp4, kdtp5, reqTglBayar, kodPeng, kodeinstitusi string, den, total uint64) string {

	queryInsertSPO := fmt.Sprintf("insert into SPO.PEMBAYARAN_SPPT (KD_PROPINSI, KD_DATI2, "+
		" KD_KECAMATAN, KD_KELURAHAN, KD_BLOK, NO_URUT, KD_JNS_OP, "+
		" THN_PAJAK_SPPT, PEMBAYARAN_SPPT_KE, KD_KANWIL_BANK, KD_KPPBB_BANK, "+
		" KD_BANK_TUNGGAL, KD_BANK_PERSEPSI, KD_TP, DENDA_SPPT, "+
		" JML_SPPT_YG_DIBAYAR, TGL_PEMBAYARAN_SPPT, TGL_REKAM_BYR_SPPT,"+
		" NIP_REKAM_BYR_SPPT, REV_FLAG, FLAG_KIRIM, PENGESAHAN, TGL_REV, "+
		" TGL_TRX, STAN, KODE_INSTITUSI) "+
		" values (35, 12, "+
		" '%s', '%s', '%s', '%s', '%v', "+
		" '%v', 1, '%v', '%v', "+
		" '%v', '%v', '%v', %v, "+
		" %v, TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS'), TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS'),"+
		" '%v', null, null, %v, null, "+
		" null, null, '%v')", kec, kel, blok, urut, kdjnsop, thn, kdtp1, kdtp2, kdtp3, kdtp4, kdtp5, den, total, reqTglBayar, reqTglBayar, kodeinstitusi[0:9], kodPeng, kodeinstitusi)

	qryUpdatePay := fmt.Sprintf("update PBB.PEMBAYARAN_SPPT "+
		"set KD_KANWIL_BANK = '%v', "+
		"KD_KPPBB_BANK = '%v', "+
		"KD_BANK_TUNGGAL = '%v', "+
		"KD_BANK_PERSEPSI = '%v', "+
		"KD_TP = '%v',"+
		"NIP_REKAM_BYR_SPPT = '%v' "+
		"where KD_KECAMATAN = '%v' and KD_KELURAHAN = '%v' and KD_BLOK = '%v' "+
		"and NO_URUT = '%v' and KD_JNS_OP = '%v' and THN_PAJAK_SPPT = '%v'", kdtp1, kdtp2, kdtp3, kdtp4, kdtp5, kodeinstitusi[0:9], kec, kel, blok, urut, kdjnsop, thn)

	//qryInsertPay := fmt.Sprintf("INSERT INTO PBB.PEMBAYARAN_SPPT (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN, "+
	//	"KD_BLOK, NO_URUT, KD_JNS_OP, THN_PAJAK_SPPT, PEMBAYARAN_SPPT_KE, KD_KANWIL_BANK, KD_KPPBB_BANK, "+
	//	"KD_BANK_TUNGGAL, KD_BANK_PERSEPSI, KD_TP, DENDA_SPPT, JML_SPPT_YG_DIBAYAR, TGL_PEMBAYARAN_SPPT, "+
	//	"NIP_REKAM_BYR_SPPT, JENIS_BAYAR, PAJAK_POOL, CREA_USER, CREA_DATE, KETERANGAN) "+
	//	"VALUES ('35', '12', '%v', '%v', "+
	//	"'%v', '%v', '%v', '%v', 1, '%v', '%v', "+
	//	"'%v', '%v', '%v', %v, %v, TO_DATE('%v', 'YYYY-MM-DD HH24:MI:SS'), "+
	//	"'%v', null, 'N', null, null, %v)", kec, kel, blok, urut, kdjnsop, thn, kdtp1, kdtp2, kdtp3, kdtp4, kdtp5, den, total, reqTglBayar, kodeinstitusi[0:9], kodPeng)

	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)

	//hasilPay, err := kon.Exec(qryInsertPay)
	//if err != nil {
	//	logrus.Fatalf("errornya di baris 215: %v\n%v\nhasilnya: %v", err.Error(), qryInsertPay, hasilPay)
	//} else {
	//	logrus.Infof("Insert pembayaran sppt Sukses")
	//}

	hasilSPO, err := kon.Exec(queryInsertSPO)
	hasilPay, err := kon.Exec(qryUpdatePay)
	if err != nil {
		logrus.Fatalf("errornya di baris 223: %v\n%v\nhasilnya: %v", err.Error(), queryInsertSPO, hasilSPO)
		logrus.Fatalf("errornya di baris 224: %v\n%v\nhasilnya: %v", err.Error(), qryUpdatePay, hasilPay)
	} else {
		logrus.Infof("Insert pembayaran %v-%v-%v-%v-%v|%v SPO Sukses\nUpdate Payment Sukses", kec, kel, blok, urut, kdjnsop, thn)
	}
	defer kon.Close()

	return queryInsertSPO
}

func ambilKodeTP(kodeinstitusi string) model.KDTP {

	var kdtp model.KDTP

	ambilTP := fmt.Sprintf("select a.kd_kanwil, a.kd_kppbb, a.kd_bank_tunggal, a.kd_bank_persepsi, a.kd_tp, a.no_rek_tp "+
		"from TEMPAT_PEMBAYARAN a where a.NO_REK_TP like '%v'", kodeinstitusi)

	config.Getenvfile()
	envUser := os.Getenv("userpbb")
	envPass := os.Getenv("password")
	envAddr := os.Getenv("addrpbb")
	envPort := os.Getenv("portpbb")
	envSN := os.Getenv("servicenamepbb")

	kon, _ := database.KonekOracle(envUser, envPass, envAddr, envPort, envSN)
	rows, err := kon.Query(ambilTP)
	if err != nil {
		logrus.Fatalf("errornya di baris 249 : %v", err.Error())
	}

	for rows.Next() {
		if err := rows.Scan(&kdtp.KD_KANWIL, &kdtp.KD_KPPBB, &kdtp.KD_BANK_TUNGGAL, &kdtp.KD_BANK_PERSEPSI, &kdtp.KD_TP, &kdtp.KODE_INSTITUSI); err != nil {
			logrus.Fatalf("errornya di baris 254 : %v", err.Error())
		}
	}
	defer rows.Close()

	//logrus.Infof("isi dari kodetp baris 258 : %v", kdtp)

	if kdtp.KD_KANWIL == "" {
		kdtp.KD_KANWIL = "12"
		kdtp.KD_KPPBB = "10"
		kdtp.KD_BANK_TUNGGAL = "00"
		kdtp.KD_BANK_PERSEPSI = "00"
		kdtp.KD_TP = "04"
		kdtp.KODE_INSTITUSI = "090909091"
	}
	defer kon.Close()

	return kdtp
}
