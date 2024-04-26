package dummyjson

type Todo struct {
	Todo      string `json:"todo,omitempty"`
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"userId,omitempty"`
	Completed bool   `json:"completed,omitempty"`
}

type TodosResp struct {
	Todos []*Todo `json:"todos,omitempty"`
	Total int     `json:"total,omitempty"`
	Skip  int     `json:"skip,omitempty"`
	Limit int     `json:"limit,omitempty"`
}
