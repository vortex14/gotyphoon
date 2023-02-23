package rod

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
)

func getDeviceOfBrowser(options Options) devices.Device {
	var device devices.Device
	if len(options.Device.UserAgent) == 0 {
		device = devices.IPadPro
	} else {
		device = options.Device
	}
	return device
}

func CreateBaseBrowser(options Options) *rod.Browser {
	
	browser := rod.New().ControlURL(CreateLauncher(options).MustLaunch())

	if options.Timeout > 0 {
		browser = browser.Timeout(options.Timeout)
	}

	return browser.DefaultDevice(getDeviceOfBrowser(options))
}