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

type TaskPurpose struct {
	Event   string `json:"event"`
	Payload struct {
		Title       string        `json:"title"`
		Timestamp   int64         `json:"timestamp"`
		ColumnId    string        `json:"columnId"`
		Description string        `json:"description"`
		Archived    bool          `json:"archived"`
		Completed   bool          `json:"completed"`
		Subtasks    []interface{} `json:"subtasks"`
		Assigned    []string      `json:"assigned"`
		CreatedBy   string        `json:"createdBy"`
		Checklists  []struct {
			Title string `json:"title"`
			Items []struct {
				IsCompleted bool   `json:"isCompleted"`
				Title       string `json:"title"`
			} `json:"items"`
		} `json:"checklists"`
		Id      string        `json:"id"`
		Parents []interface{} `json:"parents"`
	} `json:"payload"`
	PrevData struct {
		Title       string        `json:"title"`
		Timestamp   int64         `json:"timestamp"`
		ColumnId    string        `json:"columnId"`
		Description string        `json:"description"`
		Archived    bool          `json:"archived"`
		Completed   bool          `json:"completed"`
		Subtasks    []interface{} `json:"subtasks"`
		CreatedBy   string        `json:"createdBy"`
		Checklists  []struct {
			Title string `json:"title"`
			Items []struct {
				IsCompleted bool   `json:"isCompleted"`
				Title       string `json:"title"`
			} `json:"items"`
		} `json:"checklists"`
		Id      string        `json:"id"`
		Parents []interface{} `json:"parents"`
	} `json:"prevData"`
	FromUserId string `json:"fromUserId"`
}

type ResponceUser struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"isAdmin"`
	RealName     string `json:"realName"`
	Status       string `json:"status"`
	LastActivity int    `json:"lastActivity"`
}

type ResponceTask struct {
	Title       string        `json:"title"`
	Timestamp   int64         `json:"timestamp"`
	ColumnId    string        `json:"columnId"`
	Description string        `json:"description"`
	Archived    bool          `json:"archived"`
	Completed   bool          `json:"completed"`
	Subtasks    []interface{} `json:"subtasks"`
	Assigned    []string      `json:"assigned"`
	CreatedBy   string        `json:"createdBy"`
	Checklists  []struct {
		Title string `json:"title"`
		Items []struct {
			IsCompleted bool   `json:"isCompleted"`
			Title       string `json:"title"`
		} `json:"items"`
	} `json:"checklists"`
	IdTaskCommon  string `json:"idTaskCommon"`
	IdTaskProject string `json:"idTaskProject"`
	Id            string `json:"id"`
}
