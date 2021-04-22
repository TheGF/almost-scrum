// +build arm

package web

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func startSystray() {
}

func onReady() {
}

func onExit() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("listen: %s\n", err)
	}
}