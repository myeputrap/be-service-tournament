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
	slog.Info("[Repository][SaveFile] Start", slog.String("originalFilename", asset.Filename), slog.String("path", path))

	src, err := asset.Open()
	if err != nil {
		slog.Error("[Repository][SaveFile] Failed to open file", slog.Any("error", err))
		return
	}
	defer func() {
		closeErr := src.Close()
		if closeErr != nil {
			slog.Error("[Repository][SaveFile] Failed to close source", slog.Any("error", closeErr))
		}
	}()

	n := strings.LastIndexByte(asset.Filename, '.')
	if n < 0 {
		err = fmt.Errorf("invalid file extension in filename: %s", asset.Filename)
		slog.Error("[Repository][SaveFile] Extension parse failed", slog.Any("error", err))
		return
	}
	extension := asset.Filename[n:]
	fileName = fmt.Sprintf("%s%s", generatedFileName, extension)

	// Ensure directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Warn("[Repository][SaveFile] Path does not exist, attempting to create", slog.String("path", path))
		if mkErr := os.MkdirAll(path, 0755); mkErr != nil {
			slog.Error("[Repository][SaveFile] Failed to create directory", slog.Any("error", mkErr))
			return "", mkErr
		}
	}

	dst, err := os.Create(path + "/" + fileName)
	if err != nil {
		slog.Error("[Repository][SaveFile] Failed to create file", slog.Any("error", err))
		return
	}
	defer func() {
		closeErr := dst.Close()
		if closeErr != nil {
			slog.Error("[Repository][SaveFile] Failed to close destination", slog.Any("error", closeErr))
		}
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		slog.Error("[Repository][SaveFile] Failed to write file", slog.Any("error", err))
		return
	}

	slog.Info("[Repository][SaveFile] File saved successfully", slog.String("fileName", fileName))
	return
}
