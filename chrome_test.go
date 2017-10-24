package main

import (
	"testing"
)

func TestInstallChromeManifest(t *testing.T) {
	err := InstallChromeManifest()
	if err != nil {
		t.Error(err)
	}

}
