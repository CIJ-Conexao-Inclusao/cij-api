package integration

import "github.com/cloudinary/cloudinary-go/v2"

func CloudinaryConnect(url string) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromURL(url)

	if err != nil {
		panic("failed to connect cloudinary")
	}

	return cld
}