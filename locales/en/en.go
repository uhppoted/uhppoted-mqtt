package en

var Dictionary = map[string]string{
	// event type
	"event.type.0":   "none",
	"event.type.1":   "card swipe",
	"event.type.2":   "door",
	"event.type.3":   "alarm",
	"event.type.255": "overwritten",

	// event direction
	"event.direction.1": "in",
	"event.direction.2": "out",

	// event reason
	"event.reason.1":  "swipe",
	"event.reason.5":  "swipe:denied (system)",
	"event.reason.6":  "no access rights",
	"event.reason.7":  "incorrect password",
	"event.reason.8":  "anti-passback",
	"event.reason.9":  "more cards",
	"event.reason.10": "first card open",
	"event.reason.11": "door is normally closed",
	"event.reason.12": "interlock",
	"event.reason.13": "not in allowed time period",
	"event.reason.15": "invalid timezone",
	"event.reason.18": "access denied",
	"event.reason.20": "push button ok",
	"event.reason.23": "door opened",
	"event.reason.24": "door closed",
	"event.reason.25": "door opened (supervisor password)",
	"event.reason.28": "controller power on",
	"event.reason.29": "controller reset",
	"event.reason.31": "pushbutton invalid (door locked)",
	"event.reason.32": "pushbutton invalid (offline)",
	"event.reason.33": "pushbutton invalid (interlock)",
	"event.reason.34": "pushbutton invalid (threat)",
	"event.reason.37": "door open too long",
	"event.reason.38": "forced open",
	"event.reason.39": "fire",
	"event.reason.40": "forced closed",
	"event.reason.41": "theft prevention",
	"event.reason.42": "24x7 zone",
	"event.reason.43": "emergency",
	"event.reason.44": "remote open door",
	"event.reason.45": "remote open door (USB reader)",

	// doors
	"door.control.1": "normally open",
	"door.control.2": "normally closed",
	"door.control.3": "controlled",
}
