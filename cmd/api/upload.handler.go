package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/spf13/viper"
)

type UploadFileInput struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	Base64   string `json:"base64"`
}

func (app *application) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	cloudName := viper.Get("CLOUDINARY_NAME")
	cloudApiKey := viper.Get("CLOUDINARY_KEY")
	cloudSecret := viper.Get("CLOUDINARY_SECRET")
	cld, _ := cloudinary.NewFromParams(cloudName.(string), cloudApiKey.(string), cloudSecret.(string))

	input := &UploadFileInput{}

	err := json.NewDecoder(r.Body).Decode(input)

	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
		return
	}

	var ctx = context.Background()

	resp, err := cld.Upload.Upload(ctx, input.Base64, uploader.UploadParams{PublicID: input.FileName, Folder: "/healingjourney/images", ResourceType: input.FileType})

	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": err.Error()})
		res.status(http.StatusBadRequest).json(errMsg)
		return
	}

	res.status(http.StatusCreated).json(envelop{"message": "success", "url": resp.SecureURL, "width": resp.Width, "height": resp.Height})
}
