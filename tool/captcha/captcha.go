package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = base64Captcha.DefaultMemStore

type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

func DriverDigitGenerate() (id string, bs64 string, err error) {
	params := configJsonBody{}
	params.DriverString = base64Captcha.NewDriverString(80, 240, 5, 2, 4, "0123456789", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"})
	driver := params.DriverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)

	return c.Generate()
}
