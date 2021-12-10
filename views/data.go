package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
)

// Alert is ued to render bootstrap alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// data is the top levle structure that views expects data to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}
