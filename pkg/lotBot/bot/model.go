package bot

type StudentData struct {
	Tgid     string `json:"tgId"`
	Name     string `json:"Name"`
	Birthday string `json:"birthday"`
	City     string `json:"city"`
	Skill    string `json:"skill"`
	Email    string `json:"email"`
}

type BusinesData struct {
	Tgid                  string `json:"tgId"`
	CompanyName           string `json:"CompanyName"`
	INN                   string `json:"INN"`
	FieldOfActivity       string `json:"FieldOfActivity"`
	ContactPersonFullName string `json:"ContactPersonFullName"`
	ContactPersonPhone    string `json:"ContactPersonPhone"`
}

type TaskData struct {
	Tgid        string   `json:"tgId"`
	Description string   `json:"description"`
	IMG         []string `json:"IMG"`
	Link        string   `json:"Link"`
	Deadline    string   `json:"deadline"`
	SlotCall    string   `json:"slotCall"`
}
