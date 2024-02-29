package cv_models

import "encoding/json"

type PersonalInfo struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Education struct {
	Period Period `json:"period"`
	School string `json:"school"`
	Degree string `json:"degree"`
}

type Experience struct {
	Period  Period `json:"period"`
	Company string `json:"company"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type CV struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Title     string `json:"title"`
	Abstract  string `json:"abstract"`
	AvatarUrl string `json:"avatarUrl"`

	Personal   PersonalInfo `json:"personal"`
	Education  []Education  `json:"education"`
	Experience []Experience `json:"experience"`
}

func (cv CV) GetId() string {
	return cv.Id
}

func CVFromJSON(jsonData []byte) CV {
	var cv CV
	if err := json.Unmarshal(jsonData, &cv); err != nil {
		panic("unable to parse cv json")
	}
	return cv
}
