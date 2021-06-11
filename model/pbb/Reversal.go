package model

type StructReqReversal struct {
	NOP       string       `json:"Nop"`
	REFERENCE string       `json:"Reference"`
	DATETIME  string       `json:"DateTime"`
	TAGIHAN   []RevTagihan `json:"Tagihan"`
}

type StructResReversal struct {
	NOP       string    `json:"Nop"`
	REFERENCE string    `json:"Reference"`
	STATUS    RevStatus `json:"Status"`
}

type RevTagihan struct {
	TAHUN string `json:"Tahun"`
}

type RevStatus struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}

type RevStatusError struct {
	STATUSERROR RevError `json:"Status"`
}

type RevError struct {
	ISERROR      string `json:"IsError"`
	RESPONSECODE string `json:"ResponseCode"`
	ERRORDESC    string `json:"ErrorDesc"`
}
