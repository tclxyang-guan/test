package models

type Result struct {
	Status bool
	Msg    string
	Data   interface{}
}

func CreateResult(result *Result, d interface{}, msg string, err error) {
	if err != nil {
		result.Status = false
		result.Msg = msg
		result.Data = nil
	} else {
		result.Status = true
		result.Msg = "success"
		result.Data = d
	}
}
func GetResult(d interface{}, msg string, err error) (result *Result) {
	if err != nil {
		result.Status = false
		result.Msg = msg
		result.Data = nil
	} else {
		result.Status = true
		result.Msg = "success"
		result.Data = d
	}
	return
}
