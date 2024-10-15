package requests

type ReqTasks struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      string `json:"user_id,omitempty"`
}

type ReqTasksUpdate struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	Status      string `json:"status,omitempty"`
}
