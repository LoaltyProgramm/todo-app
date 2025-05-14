package models

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResp struct {
	Tasks []Task `json:"tasks"`
}

type TaskRequestPassword struct {
	Password string `json:"password"`
}