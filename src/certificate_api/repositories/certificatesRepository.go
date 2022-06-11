package repositories

import (
	"certificate_api/models"
	"fmt"
)

type CertificatesRepo struct {
	certificatesList models.CertificateModel
}

type CertificateRequestsRepo struct {
	requestsList models.CertificateRequestsModel
}

