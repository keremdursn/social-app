package database

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func ConnectToCloudinary(imageBytes []byte) (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error load .env file")
		panic(err)
	}

	cloudinaryName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Yeni bir UUID oluştur
	id := uuid.New()

	// UUID'yi string olarak al
	idString := id.String()

	cld, err := cloudinary.NewFromParams(cloudinaryName, apiKey, apiSecret)
	if err != nil {
		log.Println("cloud file could not be created.")
		panic(err)
	}

	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(imageBytes), uploader.UploadParams{PublicID: idString})
	if err != nil {
		log.Println("image failed to load.")
		panic(err)
	}

	url := GetPhoto(idString)
	log.Println(resp)
	return idString, url, nil
}

func GetPhoto(image string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error load .env file")
		panic(err)
	}

	cloudinaryName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudinaryName, apiKey, apiSecret)
	if err != nil {
		log.Println("cloud file could not be created.")
		panic(err)
	}

	var ctx = context.Background()
	resp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: image})
	if err != nil {
		log.Println("image could not get:", err)	
		panic(err)
	}

	return resp.SecureURL
}
