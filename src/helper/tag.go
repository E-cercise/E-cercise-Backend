package helper

import "strings"

var TagGroupMapping = map[string]string{
	// 💪 Muscle Tags
	"abs":       "muscle",
	"core":      "muscle",
	"arms":      "muscle",
	"shoulders": "muscle",
	"chest":     "muscle",
	"back":      "muscle",
	"legs":      "muscle",
	"glutes":    "muscle",
	"full-body": "muscle",

	// 🎯 Goal Tags
	"tone":         "goal",
	"build-muscle": "goal",
	"weight-loss":  "goal",
	"endurance":    "goal",
	"rehab":        "goal",
	"mobility":     "goal",
	"flexibility":  "goal",

	// 🧠 Experience Level Tags
	"beginner-friendly": "experience",
	"intermediate":      "experience",
	"advanced":          "experience",
	"athlete":           "experience",
	"elderly":           "experience",

	// ⚙️ Functionality / Features
	"resistance":         "feature",
	"weighted":           "feature",
	"stretching":         "feature",
	"cardio":             "feature",
	"calisthenics":       "feature",
	"pull-up":            "feature",
	"dip":                "feature",
	"ab-machine":         "feature",
	"rowing":             "feature",
	"cable":              "feature",
	"tower":              "feature",
	"barbell-compatible": "feature",
	"multi-function":     "feature",
	"gym-grade":          "feature",

	// 🏃 Usability Tags
	"bodyweight":     "usage",
	"low-impact":     "usage",
	"joint-friendly": "usage",
	"post-injury":    "usage",
	"compact":        "usage",
	"adjustable":     "usage",
	"budget":         "usage",
	"foldable":       "usage",
	"portable":       "usage",
}

func GetTagGroup(tagName string) string {
	group := TagGroupMapping[strings.ToLower(tagName)]
	if group == "" {
		group = "general"
	}

	return group
}
