package handlers

var flash Flash

type Flash struct {
	Type    string
	Message string
}

func (f *Flash) Set(_type, msg string) {
	f.Type = _type
	f.Message = msg
}
func (f *Flash) Delete() {
	f.Type = ""
	f.Message = ""
}
