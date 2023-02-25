package rod

import (
	"strings"

	. "github.com/go-rod/rod/lib/launcher"
)

const (
	MuteAudio            = "mute-audio"
	DisableGPU           = "disable-gpu"
	HideScrollBars       = "hide-scrollbars"
	UseMockKeychain      = "use-mock-keychain"
	DisableNotifications = "disable-notifications"
	DisableCrashReporter = "disable-crash-reporter"
)

func createLauncher() *Launcher {
	return New().
		Delete(UseMockKeychain).
		Set(MuteAudio, "true").
		Set(DisableGPU, "true").
		Set(HideScrollBars, "true").
		Set(DisableCrashReporter, "true").
		Set(DisableNotifications, "true")
}

func CreateLauncher(options Options) *Launcher {
	var _launcher *Launcher
	if options.Debug {
		_launcher = createLauncher().Headless(false)
	} else {
		_launcher = createLauncher()
	}

	if len(options.Proxy) > 0 {
		_proxy := strings.Split(options.Proxy, "@")
		if len(_proxy) == 2 {
			_launcher.Proxy(_proxy[1])
		} else {
			_launcher.Proxy(options.Proxy)
		}

	}

	return _launcher
}
