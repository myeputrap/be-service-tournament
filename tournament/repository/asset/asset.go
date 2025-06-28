package asset

import (
	"be-service-tournament/constant"
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"fmt"

	"io"
	"mime/multipart"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
)

type assetRepository struct{}

func NewAssetRepository() domain.AssetRepository {
	return &assetRepository{}
}

func (r *assetRepository) SaveAsset(req domain.RequestPaymentProffImage) (filename string, status int, err error) {
	assetpath := viper.GetString("server.http.asset_path")
	slog.Debug(req.Images.Filename)
	if req.Images.Size > constant.MaxSizeBanner {
		slog.Error("[Handler][SaveAsset-PostFeatureImage]The image must be less than 10Mb", "Err", nil)
		status = domain.StatusInternalServerError
		return
	}

	err = helper.CheckFileImage(req.Images)
	if err != nil {
		slog.Error("[Handler][SaveAsset-PostFeatureImage]check file", "Err", err.Error())
		status = domain.StatusInternalServerError
		return
	}

	filename, err = helper.Save(req.Images, assetpath+"/images/features", "")
	if err != nil {
		slog.Error("[Handler][SaveAsset-PostFeatureImage]save file", "Err", err.Error())
		status = domain.StatusInternalServerError
		return
	}
	return
}

func (r *assetRepository) Remove(path string, generatedFileName string) (err error) {
	param := path + "/" + generatedFileName
	err = os.Remove(param)
	if err != nil {
		slog.Error("[Repository][DeleteFile]", "Err", err.Error())
		return
	}
	return
}

func (r *assetRepository) SaveFile(asset *multipart.FileHeader, path string, generatedFileName string) (fileName string, err error) {
	src, err := asset.Open()
	if err != nil {
		return
	}
	defer func() {
		err = src.Close()
	}()

	n := strings.LastIndexByte(asset.Filename, '.')
	extension := asset.Filename[n:]
	fileName = fmt.Sprintf("%s%s", generatedFileName, extension)
	dst, err := os.Create(path + "/" + fileName)
	if err != nil {
		slog.Error("[Repository][SaveFile]", "Err", err.Error())
		return
	}
	defer func() {
		err = dst.Close()
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		slog.Error("[Repository][SaveFile]", "Err", err.Error())
		return
	}
	return
}
