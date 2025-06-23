package helper

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type UserLogin struct {
	ID  int32  `json:"id"`
	HID *int32 `json:"hid"`
	RID *int32 `json:"rid"`
}

func GetUserLogin(c context.Context) (user UserLogin) {
	user = c.Value("user_login").(UserLogin)
	return
}

func ContainsAny(tags []string, targets []string) bool {
	for _, tag := range tags {
		for _, target := range targets {
			if tag == target {
				return true
			}
		}
	}
	return false
}

func CreateLanguageImageUrl(languageCode string) string {
	return fmt.Sprintf("%s/%s/%s/%s.svg", viper.GetString("server.http.base_url"), "images", "languages", languageCode)
}

func CreateHotelLogoImageUrl(hotelLogo string) string {
	return fmt.Sprintf("%s/%s/%s/%s", viper.GetString("server.http.base_url"), "images", "hotels", hotelLogo)
}

func CreateHotelBackgroundUrl(hotelBackground string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", viper.GetString("server.http.base_url"), "images", "hotels", "backgrounds", hotelBackground)
}

func CreateThemeVariantUrl(variant string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s.svg", viper.GetString("server.http.base_url"), "images", "hotels", "variants", variant)
}

func CreateAppAndroidUrl(appName string) string {
	return fmt.Sprintf("%s/%s/%s/%s", viper.GetString("server.http.base_url"), "apps", "android", appName)
}

func CreateAppAndroidVersionUrl(packageName, version string) string {
	return fmt.Sprintf("%s/%s/%s/%s.apk", viper.GetString("server.http.base_url"), "apps", "android", packageName+"_"+version)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CheckVariantEnum(variant *string) bool {
	if variant == nil {
		return true
	}

	switch *variant {
	case "flower", "circle", "quantum":
		return true
	}
	return false
}

func CheckFontTheme(font *string) bool {
	if font == nil {
		return true
	}

	switch *font {
	case "verdana", "calibri", "inter", "helvetica", "montserrat":
		return true
	}
	return false
}

func CheckHexcolorString(str *string) bool {
	if str == nil {
		return true
	}

	hexColorPattern := `^#([a-fA-F0-9]{3}|[a-fA-F0-9]{6}|[a-fA-F0-9]{8})$`
	match, _ := regexp.MatchString(hexColorPattern, *str)
	return match
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func StrPtr(s string) *string {
	return &s
}

func StringToNilIfEmpty(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

func Float64ToNilIfZero(s *float64) *float64 {
	if s != nil && *s == 0 {
		return nil
	}
	return s
}

func Int32ToNilIfZero(s *int32) *int32 {
	if s != nil && *s == 0 {
		return nil
	}
	return s
}

func IntToNilIfZero(s *int) *int {
	if s != nil && *s == 0 {
		return nil
	}
	return s
}

func Int64ToNilIfZero(s *int64) *int64 {
	if s != nil && *s == 0 {
		return nil
	}
	return s
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ParseTime(dateStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
	if err != nil {
		slog.Error("Failed to parse date", "Err", err.Error())
		return time.Time{}
	}
	return parsedTime
}

func IsISO8601(dateStr *string) bool {
	if dateStr == nil || *dateStr == "" {
		return false
	}

	layouts := []string{
		time.RFC3339,               // "2006-01-02T15:04:05Z"
		"2006-01-02T15:04:05.000Z", // With milliseconds
		"2006-01-02",               // Date only
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, *dateStr); err == nil {
			return true
		}
	}

	return false
}
func StringPtr(s string) *string {
	return &s
}

func DecodeBase64(req string) (res string, err error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		return
	}

	res = string(decodedBytes)
	return
}

func MoveFile(sourcePath, destPath string) error {
	// inputFile, err := os.Open(sourcePath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't open source file: %v", err)
	// }
	// defer inputFile.Close()

	// outputFile, err := os.Create(destPath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't open dest file: %v", err)
	// }
	// defer outputFile.Close()

	// _, err = io.Copy(outputFile, inputFile)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	// }

	// inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	// err = os.Remove(sourcePath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't remove source file: %v", err)
	// }
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		fmt.Println("Error moving file:", err)
		return err
	}
	return nil
}

func DeleteFile(sourcePath string) error {
	// inputFile, err := os.Open(sourcePath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't open source file: %v", err)
	// }
	// defer inputFile.Close()

	// outputFile, err := os.Create(destPath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't open dest file: %v", err)
	// }
	// defer outputFile.Close()

	// _, err = io.Copy(outputFile, inputFile)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	// }

	// inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	// err = os.Remove(sourcePath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't remove source file: %v", err)
	// }
	err := os.Remove(sourcePath)
	if err != nil {
		fmt.Println("Error moving file:", err)
		return err
	}
	return nil
}

func GetFileName(url string) string {
	parts := strings.Split(url, "features/")
	if len(parts) > 1 {
		return parts[1]
	}
	return "" // or handle the error if needed
}
