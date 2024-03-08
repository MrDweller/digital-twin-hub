package models

import "github.com/google/uuid"

type CertificateModel struct {
	CertificateId uuid.UUID `json:"certificateId"`
	CertFilePath  string    `json:"certFilePath"`
	KeyFilePath   string    `json:"keyFilePath"`
}
