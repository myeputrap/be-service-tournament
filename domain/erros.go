package domain

import (
	"errors"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
)

var (
	StatusSuccess                        = int(20000)
	StatusSuccessLogin                   = int(20001)
	StatusSuccessRegister                = int(20002)
	StatusSuccessLoginUnlinked           = int(20003)
	StatusSuccessCreate                  = int(20100)
	StatusBadRequest                     = int(40000)
	StatusMissingParameter               = int(40001)
	StatusNotRecognized                  = int(40002)
	StatusWrongValue                     = int(40003)
	StatusChannelCategoryViolation       = int(40004)
	StatusRegisterDuplicateEmail         = int(40005)
	StatusInvalidEmail                   = int(40006)
	StatusRoomCustomNotFound             = int(40007)
	StatusFeatureIDIsUsed                = int(40008)
	StatusGreetingIDIsExist              = int(40009)
	StatusGreetingMsgExist               = int(40010)
	StatusAppVersionExist                = int(40011)
	StatusDuplicateDataLauncher          = int(40012)
	StatusLMSPartnerIDNotValid           = int(40013)
	StatusLMSPartnerIDAlreadyUsed        = int(40014)
	StatusLMSPartnerIDStillUsedByRoom    = int(40015)
	StatusWrongFormatDate                = int(40050)
	StatusUnauthorized                   = int(40100)
	StatusUnregistered                   = int(40101)
	StatusUnlinked                       = int(40102)
	StatusInvalidEmailPassword           = int(40103)
	StatusUnauthorizedUnverified         = int(40104)
	StatusUnauthorizedBlockedAccount     = int(40105)
	StatusUnauthorizedDisabledDevice     = int(40106)
	StatusForbidden                      = int(40300)
	StatusForbiddenWrongHotelID          = int(40301)
	StatusNotFound                       = int(40400)
	StatusLanguageNotFound               = int(40401)
	StatusLanguageGreetingNotFound       = int(40402)
	StatusGeneralGreetingNotFound        = int(40403)
	StatusGeneralPersonalizationNotFound = int(40404)
	StatusGeneralMessageNotFound         = int(40405)
	StatusGetSelectLanguageNotFound      = int(40406)
	StatusGetPersonalizationNotFound     = int(40407)
	StatusGetWelcomeDeviceNotFound       = int(40408)
	StatusGetGreetingMessageNotFound     = int(40409)
	StatusLanguageNotExistInSystem       = int(40410)
	StatusRoomNotFound                   = int(40411)
	StatusAppVersionNotFound             = int(40440)
	StatusLmsLicenceNotFound             = int(40412)
	StatusTemplateNotValid               = int(40414)
	StatusInternalServerError            = int(50000)
	StatusLMSError                       = int(50001)
	StatusEmptyAccountLMS                = int(50002)
	StatusLicenseNotAvailable            = int(50003)
	StatusAccountLMSNotFound             = int(50004)
	StatusAccountLMSStillPaired          = int(50005)
	StatusRoomAccountLMSIsNilAndStatus   = int(40016)
	StatusAccountLMSIsNil                = int(40018)
	StatusAccountIsNotPaired             = int(40019)
	StatusFailSendNotifMessageGuest      = int(40020)
	StatusGuestNotFound                  = int(40413)
	StatusFailDeleteRoom                 = int(40021)
	StatusMissingNewValue                = int(40099)
	StatusPinnedValueExceeded            = int(40110)
	StatusCantUpdateVideos               = int(40111)
	StatusPlayerAlreadyRegistered        = int(40112)
	StatusWrongReferalCode               = int(40113)
)

var (
	ErrBadRequest             = errors.New("bad request")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrNotFound               = errors.New("not found")
	ErrInternalServerError    = errors.New("internal server error")
	ErrInternalServerErrorNFC = errors.New("internal server error nfc")
)

func GetHttpStatusCode(status int) int {
	switch status {
	case StatusSuccess, StatusSuccessLogin, StatusSuccessRegister, StatusSuccessLoginUnlinked:
		return fiber.StatusOK
	case StatusSuccessCreate:
		return fiber.StatusCreated
	case StatusBadRequest, StatusWrongReferalCode, StatusPlayerAlreadyRegistered, StatusCantUpdateVideos, StatusPinnedValueExceeded, StatusTemplateNotValid, StatusWrongFormatDate, StatusMissingNewValue, StatusMissingParameter, StatusAccountIsNotPaired, StatusRoomAccountLMSIsNilAndStatus, StatusNotRecognized, StatusWrongValue, StatusChannelCategoryViolation, StatusRegisterDuplicateEmail, StatusInvalidEmail, StatusRoomCustomNotFound, StatusLanguageNotExistInSystem, StatusFeatureIDIsUsed, StatusGreetingIDIsExist, StatusGreetingMsgExist, StatusAppVersionExist, StatusDuplicateDataLauncher, StatusLMSPartnerIDNotValid, StatusLMSPartnerIDAlreadyUsed, StatusLMSPartnerIDStillUsedByRoom, StatusFailSendNotifMessageGuest, StatusFailDeleteRoom:
		return fiber.StatusBadRequest
	case StatusUnauthorized, StatusUnregistered, StatusUnlinked, StatusInvalidEmailPassword, StatusUnauthorizedUnverified, StatusUnauthorizedBlockedAccount, StatusUnauthorizedDisabledDevice:
		return fiber.StatusUnauthorized
	case StatusForbidden, StatusForbiddenWrongHotelID:
		return fiber.StatusForbidden
	case StatusNotFound, StatusLanguageNotFound, StatusGuestNotFound, StatusRoomNotFound, StatusLanguageGreetingNotFound, StatusGeneralGreetingNotFound, StatusGeneralPersonalizationNotFound, StatusGeneralMessageNotFound, StatusGetSelectLanguageNotFound, StatusGetPersonalizationNotFound, StatusGetWelcomeDeviceNotFound, StatusGetGreetingMessageNotFound, StatusAppVersionNotFound, StatusLmsLicenceNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

func GetStatusGRPCErr(cd codes.Code) int {
	switch cd {
	case codes.OK:
		return StatusSuccess
	case codes.InvalidArgument:
		return StatusBadRequest
	case codes.PermissionDenied:
		return StatusForbidden
	case codes.Unauthenticated:
		return StatusUnauthorized
	case codes.NotFound:
		return StatusNotFound
	default:
		return StatusInternalServerError
	}
}

func GetCustomStatusMessage(status int, m string) string {
	switch status {
	case StatusSuccess:
		return "Success"
	case StatusSuccessCreate:
		return "Succes create"
	case StatusBadRequest:
		return "Bad request"
	case StatusMissingParameter:
		return "Missing required parameter " + m
	case StatusWrongValue:
		return "Wrong value for parameter " + m
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not found"
	case StatusInvalidEmailPassword:
		return "Username or password is wrong"
	case StatusPlayerAlreadyRegistered:
		return "Player already registered with other user"
	case StatusWrongReferalCode:
		return "Please input the correct referalCode"
	default:
		return "Internal server error"
	}
}

func GetStatusCode(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func GetStatusMessage(err error) string {
	switch err {
	case ErrBadRequest:
		return "Bad request"
	case ErrUnauthorized:
		return "Unauthorized"
	case ErrForbidden:
		return "Forbidden"
	case ErrNotFound:
		return "Not found"
	default:
		return "Internal server error"
	}
}
