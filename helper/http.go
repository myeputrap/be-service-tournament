package helper

import (
	"be-service-tournament/constant"
	"be-service-tournament/domain"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/avast/apkparser"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	svg "github.com/h2non/go-is-svg"
	// "github.com/gosimple/slug"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Response struct {
	Status   Status      `json:"status"`
	MetaData interface{} `json:"meta,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HttpResponse(c *fiber.Ctx, httpStatus int, message string, data interface{}) error {
	return c.Status(httpStatus).JSON(&response{
		Code:    httpStatus,
		Message: message,
		Data:    data,
	})
}

const (
	DateFormat    = "2006-01-02"
	RFC3339Format = time.RFC3339
)

func ParseDate(dateStr string) (time.Time, error) {
	// Try parsing in RFC3339 format
	if date, err := time.Parse(RFC3339Format, dateStr); err == nil {
		return date, nil
	}
	// Try parsing in DateFormat (yyyy-mm-dd)
	if date, err := time.Parse(DateFormat, dateStr); err == nil {
		return date, nil
	}
	return time.Time{}, errors.New("invalid date format")
}

func NewMetaData(totalData int32, page int32, limit int32, sort string, order string) (metaData domain.MetaData) {
	totalPage := totalData / limit
	if totalData%limit > 0 {
		totalPage++
	}
	metaData.TotalData = uint(totalData)
	metaData.TotalPage = uint(totalPage)
	metaData.Page = uint(page)
	metaData.Limit = uint(limit)
	metaData.Sort = sort
	metaData.Order = order
	return
}

func NewResponse(status int, message string, metaData interface{}, data interface{}) (response Response) {
	response.Status.Code = status
	response.Status.Message = message
	response.MetaData = metaData
	response.Data = data
	return
}
func SaveFile(asset *multipart.FileHeader, prefixFileName string, folderName string, req string) (fileName string, err error) {
	src, err := asset.Open()
	if err != nil {
		return
	}
	defer func() {
		err = src.Close()
	}()
	n := strings.LastIndexByte(asset.Filename, '.')
	// fn := asset.Filename[:n]
	extension := asset.Filename[n:]
	// fn = slug.Make(fn)
	fileName = fmt.Sprintf("%s%s", req, extension)
	dst, err := os.Create(fileName)
	if err != nil {
		slog.Error("[Helper][SaveFile]", "Err", err.Error())
		return
	}
	defer func() {
		err = dst.Close()
	}()
	_, err = io.Copy(dst, src)
	if err != nil {
		slog.Error("[Helper][SaveFile]", "Err", err.Error())
		return
	}
	return
}
func Save(asset *multipart.FileHeader, path string, generatedFileName string) (fileName string, err error) {
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
		slog.Error("[Helper][SaveFile]", "Err", err.Error())
		return
	}
	defer func() {
		err = dst.Close()
	}()
	_, err = io.Copy(dst, src)
	if err != nil {
		slog.Error("[Helper][SaveFile]", "Err", err.Error())
		return
	}
	return
}
func CheckFileImage(asset *multipart.FileHeader) (err error) {
	mimeType := asset.Header["Content-Type"][0]
	if mimeType != constant.MimeTypeImageJpeg && mimeType != constant.MimeTypeImagePng && mimeType != constant.MimeTypeImageWebp {
		err = errors.New("bad picture format, should be image/jpeg or image/png or image/webp")
	}
	return
}

func CheckFolderOrFile(path string) (isExist bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		isExist = true
		return
	}
	if os.IsNotExist(err) {
		err = nil
		return
	}
	return
}

func CheckFileImageJpegOrPng(asset *multipart.FileHeader) (err error) {
	// Open the file
	fileContent, err := asset.Open()
	if err != nil {
		return err
	}
	defer func() {
		err = fileContent.Close()
	}()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return err
	}

	// Check if the file is an image
	if !filetype.IsImage(content) {
		err = fmt.Errorf("%w: file type is not image", domain.ErrBadRequest)
		return err
	}

	kind, err := filetype.Match(content)
	if err != nil {
		return err
	}

	if kind != types.Get("jpg") && kind != types.Get("png") {
		err = fmt.Errorf("%w: bad picture format, should be image/jpeg or image/png", domain.ErrBadRequest)
		return err
	}

	return
}

func IsSVG(asset *multipart.FileHeader) (err error) {
	fileContent, err := asset.Open()
	if err != nil {
		slog.Error("[Helper][IsSVG] ", "Err", err.Error())
		return
	}
	defer func() {
		err = fileContent.Close()
	}()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		slog.Error("[Helper][IsSVG] ", "Err", err.Error())
		return
	}

	if !svg.Is(content) {
		err = errors.New("file should be svg")
		slog.Error("[Helper][IsSVG] ", "Err", err.Error())
		return
	}
	return
}

func CheckFile(asset *multipart.FileHeader) (err error) {
	// Open the file
	fileContent, err := asset.Open()
	if err != nil {
		return err
	}
	defer func() {
		err = fileContent.Close()
	}()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return err
	}

	// Check if the file is video or image
	if !filetype.IsVideo(content) && !filetype.IsImage(content) {
		err = fmt.Errorf("%w: file type is not image or video", domain.ErrBadRequest)
		return err
	}

	kind, err := filetype.Match(content)
	if err != nil {
		return err
	}

	if kind != types.Get("jpg") && kind != types.Get("png") && kind != types.Get("mp4") && kind != types.Get("mov") && kind != types.Get("avi") {
		err = fmt.Errorf("%w: bad picture format, should be image/jpeg or image/png or video/mp4 or video/quicktime or video/x-msvideo", domain.ErrBadRequest)
		return err
	}

	return
}
func CheckFileTusD(content []byte) (extension string, err error) {
	// Check if the file is video or image
	if !filetype.IsVideo(content) && !filetype.IsImage(content) {
		err = fmt.Errorf("%w: file type is not image or video", domain.ErrBadRequest)
		return
	}

	kind, err := filetype.Match(content)
	if err != nil {
		return
	}

	if kind != types.Get("mp4") && kind != types.Get("mov") && kind != types.Get("avi") && kind != types.Get("mkv") {
		err = fmt.Errorf("%w: bad picture format, should be image/jpeg or image/png or video/mp4 or video/quicktime or video/x-msvideo", domain.ErrBadRequest)
		return
	}

	return kind.Extension, nil
}

func ParseTIMERFC3339(alarmStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, alarmStr)
	if err != nil {
		return time.Time{}, errors.New("wrong date format for alarm time")
	}
	return parsedTime, nil
}

func CheckFileAndroidApp(asset *multipart.FileHeader) (err error) {
	if filepath.Ext(asset.Filename) != ".apk" {
		err = fmt.Errorf("%w: app format is not .apk", domain.ErrBadRequest)
		return
	}

	// Open the file
	fileContent, err := asset.Open()
	if err != nil {
		return err
	}
	defer func() {
		err = fileContent.Close()
	}()

	// Using bytes.Buffer as writer because doesn't need the android manifest to be printed out
	// Use os.Stdout as writer instead if printing the android manifest to the console is required
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	// enc.Indent("", "\t")

	zipErr, resErr, manErr := apkparser.ParseApkReader(fileContent, enc)
	if zipErr != nil {
		err = fmt.Errorf("%w: failed to open the APK: %s", domain.ErrBadRequest, zipErr.Error())
		return
	}

	if resErr != nil {
		err = fmt.Errorf("%w: failed to parse resources: %s", domain.ErrBadRequest, resErr.Error())

	}
	if manErr != nil {
		if err != nil {
			err = fmt.Errorf("%s; additionally, %w: failed to parse AndroidManifest.xml: %s",
				err.Error(), domain.ErrBadRequest, manErr.Error())
		} else {
			err = fmt.Errorf("%w: failed to parse AndroidManifest.xml: %s", domain.ErrBadRequest, manErr.Error())
		}
		return
	}

	return
}

func FindFirstMatch(s string, words map[string]string) string {
	for key := range words { // Iterate over the keys of the map
		if strings.Contains(s, key) {
			return words[key] // Return the corresponding value from the map
		}
	}
	return "" // No match found
}
