package chrome

import (
	"os"
)

// Chrome contains information about a Google Chrome
// instance, with methods to run on it.
type Chrome struct {
	Resolution       string
	ChromeTimeout    int
	ChromeTimeBudget int
	Path             string
	UserAgent        string
	Argvs            []string

	ScreenshotPath string
}

// Setup configures a Chrome struct with the path
// specified to what is available on this system.
func (chrome *Chrome) Setup() {

	chrome.chromeLocator()
}

// ChromeLocator looks for an installation of Google Chrome
// and returns the path to where the installation was found
func (chrome *Chrome) chromeLocator() {

	// if we already have a path to chrome (say from a cli flag),
	// check that it exists. If not, continue with the finder logic.
	_, err := os.Stat(chrome.Path);
	os.IsNotExist(err)

	// Possible paths for Google Chrome or chromium to be at.
	paths := []string{
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}

	for _, path := range paths {

		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		chrome.Path = path
	}

	// final check to ensure we actually found chrome
	if chrome.Path == "" {
	}
}
