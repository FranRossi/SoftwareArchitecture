package repositories

import (
	"certificate_api/models"
)


func StoreRequest(cert *models.CertificateModel) error {
	client := connections.GetInstanceMongoClient()
	certificatesDatabase := client.Database("certificates")
	requestsCollection := certificatesDatabase.Collection("certificates")
	_, err2 := requestsCollection.InsertOne(context.TODO(), cert)
	if err2 != nil {
		fmt.Println("error storing certificate")
		if err2 == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err2)
	}
	return nil
}
