package model

type StructReqInquiry struct {
	NOP       string `json:"Nop"`
	MASAPAJAK string `json:"MasaPajak"`
	// DATETIME      string `json:"DateTime" sql:"DEFAULT:current_timestamp"` // coba hilangkan warning
	DATETIME      string `json:"DateTime"` // coba hilangkan warning
	MERCHANT      string `json:"Merchant"`
	KODEINSTITUSI string `json:"KodeInstitusi"`
	NOHP          string `json:"NoHp"`
	EMAIL         string `json:"Email"`
}

type StructResInquiry struct {
	NOP           string       `json:"Nop"`
	NAMA          string       `json:"Nama"`
	KELURAHAN     string       `json:"Kelurahan"`
	KODEKP        string       `json:"KodeKp"`
	KODEINSTITUSI string       `json:"KodeInstitusi"`
	NOHP          string       `json:"NoHp"`
	EMAIL         string       `json:"Email"`
	TAGIHAN       []InqTagihan `json:"Tagihan"`
	STATUS        InqStatus    `json:"Status"`
}

type StructResPelayanan struct {
	NAMA      string `json:"Nama"`
	TAHUN     string `json:"Tahun"`
	POKOK     string `json:"Pokok"`
	DENDA     string `json:"Denda"`
	TOTAL     string `json:"Total"`
	LUASINDUK string `json:"LuasInduk"`
	KODE      string `json:"Kode"`
	ALAMATWP  string `json:"AlamatWP"`
	STATUS    string `json:"Status"`
	LUNAS    string `json:"Lunas"`
	NJOPBUMI    string `json:"NjopBumi"`
	NJOPBNG    string `json:"NjopBng"`
	TANGGALBAYAR    string `json:"TanggalBayar"`
}

type InqTagihan struct {
	TAHUN string `json:"Tahun"`
	POKOK uint64 `json:"Pokok"`
	// DENDA sql.NullInt32 `json:"Denda"`
	DENDA uint64 `json:"Denda"`
	TOTAL uint64 `json:"Total"`
	//LUNAS byte   `json:"Lunas"`
	//JATUHTEMPO string `json:"JatuhTempo"`
}

type InqNama struct {
	NAMA      string `json:"Nama"`
	KELURAHAN string `json:"Kelurahan"`
}

type InqStatus struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}

type InqStatusError struct {
	STATUSERROR InqError `json:"Status"`
}

type InqError struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}
