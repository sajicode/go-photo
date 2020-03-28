package views

const (
	// AlertLvlError danger messages
	AlertLvlError = "danger"
	// AlertLvlWarning warning messages
	AlertLvlWarning = "warning"
	// AlertLvlInfo info messages
	AlertLvlInfo = "info"
	// AlertLvlSuccess success messages
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error
	// is encountered by our backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert struct holds bootstrap alert fields
type Alert struct {
	Level   string
	Message string
}

// Data struct encompasses bootstrap alert and extra info
type Data struct {
	Alert *Alert // * by using a pointer, Alert can be nil
	Yield interface{}
}
