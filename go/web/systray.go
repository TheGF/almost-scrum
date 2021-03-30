package web

import (
	"almost-scrum/assets"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"os"
	"runtime"
)

func startSystray() {

	switch runtime.GOOS {
		case "windows":
		default: // "linux", "freebsd", "openbsd", "netbsd"
			display, present := os.LookupEnv("DISPLAY")
			if !present || display == "" {
				logrus.Warnf("No XServer found. System Icon will not be present")
				return
			}
		}

	systray.Run(onReady, onExit)
}

func onReady() {

	var iconExt string
	if runtime.GOOS == "windows" {
		iconExt = "ico"
	} else {
		iconExt = "png"
	}
	if iconData, err := assets.Asset(fmt.Sprintf("assets/icons/grapes.%s", iconExt)); err == nil {
		systray.SetIcon(iconData)
	}
	systray.SetTitle("Almost Scrum")
	systray.SetTooltip("Scrum on the go")
	openMenu := systray.AddMenuItem("New Window", "Open a new browser page")
	autoExitMenu := systray.AddMenuItemCheckbox("Auto-exit", "Exit when you close the web page", autoExit)
	systray.AddSeparator()
	quitMenu := systray.AddMenuItem("Quit", "Quit the whole app")


	go func() {
		for {
			select {
			case <- quitMenu.ClickedCh: {
				systray.Quit()
				onExit()
			}
			case <- openMenu.ClickedCh: {
				open.Start(ashUrl)
			}
			case <- autoExitMenu.ClickedCh: {
				if autoExit {
					autoExitMenu.Uncheck()
					autoExit = false
				} else {
					autoExitMenu.Check()
					autoExit = true
				}
			}
			}
		}

	}()

}

func onExit() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("listen: %s\n", err)
	}
}