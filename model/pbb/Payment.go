package model

type StructReqPayment struct {
	NOP           string       `json:"Nop"`
	MERCHANT      string       `json:"Merchant"`
	DATETIME      string       `json:"DateTime"`
	REFERENCE     string       `json:"Reference"`
	TOTALBAYAR    uint64       `json:"TotalBayar"`
	KODEINSTITUSI string       `json:"KodeInstitusi"`
	NOHP          string       `json:"NoHp"`
	EMAIL         string       `json:"Email"`
	TAGIHAN       []PayTagihan `json:"Tagihan"`
}

type StructResPayment struct {
	NOP            string    `json:"Nop"`
	KODEPENGESAHAN string    `json:"KodePengesahan"`
	KODEKP         string    `json:"KodeKp"`
	STATUS         PayStatus `json:"Status"`
}

type PayTagihan struct {
	TAHUN string `json:"Tahun"`
}

type PayStatus struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}

type PayStatusError struct {
	STATUSERROR PayError `json:"Status"`
}

type PayError struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}

type KDTP struct {
	KD_KANWIL        string
	KD_KPPBB         string
	KD_BANK_TUNGGAL  string
	KD_BANK_PERSEPSI string
	KD_TP            string
	KODE_INSTITUSI   string
}
