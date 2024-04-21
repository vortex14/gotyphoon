package rod

import (
	"context"
	uR "github.com/EDDYCJY/fake-useragent"
	"github.com/vortex14/gotyphoon/log"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
)

var (
	D = []devices.Device{
		devices.IPad,
		devices.IPadPro,
		devices.LaptopWithHiDPIScreen,
		devices.LaptopWithMDPIScreen,
		devices.Nexus10,
		devices.KindleFireHDX,
	}
)

func getDeviceOfBrowser(options Options) devices.Device {
	var device devices.Device
	if len(options.Device.UserAgent) == 0 {
		device = devices.IPadPro
	} else if options.Device.UserAgent == "random" {
		rand.Seed(time.Now().UnixNano())
		device = D[rand.Intn(len(D))]
	} else {
		device = options.Device
	}

	if options.Screen.DevicePixelRatio != 0 {
		device.Screen = options.Screen
	}

	return device

}

func CreateBaseBrowser(context context.Context, options Options) (context.Context, *rod.Browser) {

	browser := rod.New().ControlURL(CreateLauncher(options).MustLaunch())

	if options.Debug {
		browser.Trace(true)
		browser.SlowMotion(1 * time.Second)
	}
	if options.Timeout > 0 {
		browser = browser.Timeout(options.Timeout)
	}

	device := getDeviceOfBrowser(options)

	logMap := map[string]interface{}{
		"device":     device.Title,
		"user-agent": device.UserAgent,
	}

	if options.RandomAgent {
		device.UserAgent = uR.Random()
		logMap["random_agent"] = true
	}

	context = log.PatchCtx(context, logMap)

	return context, browser.DefaultDevice(device)
}
