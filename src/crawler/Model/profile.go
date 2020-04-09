package Model

import "encoding/json"

type Profile struct {
	Name       string
	Height     int
	Age        int
	Income     string
	Marriage   string
	Occupation string
	Birthplace string
	BasicInfo  string
	DetailInfo string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var p Profile
	s, err := json.Marshal(o)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(s, &p)
	return p, nil
}
