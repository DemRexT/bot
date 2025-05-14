package bot

type StudentData struct {
	Tgid     string   `json:"tgId"`
	Name     string   `json:"Name"`
	Birthday string   `json:"birthday"`
	City     string   `json:"city"`
	Skill    []string `json:"skill"`
	Email    string   `json:"email"`
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
	TgId        string `json:"tgId"`
	NameTask    string `json:"nameTask"`
	Description string `json:"description"`
	Budget      string `json:"budget"`
	Direction   string `json:"direction"`
	Link        string `json:"Link"`
	Deadline    string `json:"deadline"`
	SlotCall    string `json:"slotCall"`
}
