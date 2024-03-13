package models

type CertificateModel struct {
	CertificateId string `json:"certificateId"`
	CertFilePath  string `json:"certFilePath"`
	KeyFilePath   string `json:"keyFilePath"`
}
