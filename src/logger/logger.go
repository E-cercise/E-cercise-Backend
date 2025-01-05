package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	// Set the log level
	Log.SetLevel(logrus.InfoLevel)

	// Set the log format to JSON (optional, you can also use TextFormatter)
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,                  // Enable colors in the terminal
		FullTimestamp:   true,                  // Include timestamps
		TimestampFormat: "2006-01-02 15:04:05", // Set timestamp format
	})

	// You can add hooks to output logs to different destinations (e.g., files)
	// Log.AddHook(myCustomHook) // Optional
}
