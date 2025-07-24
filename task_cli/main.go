package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yakumo-zi/gtool/pkg/slice"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type Status string

const Done Status = "done"
const Todo Status = "todo"
const InProgress Status = "in_progress"

const FilePath = "/tmp/todo_cli/todo.json"

type Task struct {
	Id          int       `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      Status    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%-4d\t%-24s\t%-8s\t%-10s\t%-10s", t.Id, t.Description, t.Status, t.CreatedAt.Format("2006/01/02"), t.UpdatedAt.Format("2006/01/02"))
}

type TaskContainer struct {
	Tasks  []*Task `json:"tasks,omitempty"`
	LastId int     `json:"last_id,omitempty"`
}

func LoadFromFile(filePath string) (TaskContainer, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(path.Dir(filePath), os.ModePerm)
		os.Create(filePath)
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return TaskContainer{
			Tasks: []*Task{},
		}, err
	}
	if len(bytes) == 0 {
		return TaskContainer{
			Tasks: []*Task{},
		}, nil
	}
	var tasks TaskContainer
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		return TaskContainer{
			Tasks: []*Task{},
		}, err
	}
	return tasks, nil
}

func (t *TaskContainer) Add(desc string) {
	t.Tasks = append(t.Tasks, &Task{
		Id:          t.LastId + 1,
		Description: desc,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	t.LastId += 1
}

func (t *TaskContainer) Delete(id int) error {
	idx := 0
	exists := false
	for i, task := range t.Tasks {
		if task.Id == id {
			idx = i
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("task %d is not exist", id)
	}
	t.Tasks = append(t.Tasks[:idx], t.Tasks[idx+1:]...)
	return nil
}

func (t *TaskContainer) Update(id int, desc string) error {
	updated := false
	for _, tt := range t.Tasks {
		if tt.Id == id {
			tt.Description = desc
			tt.UpdatedAt = time.Now()
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("task %d is not exist", id)
	}
	return nil
}

func (t *TaskContainer) Mark(id int, status Status) {
	t.Tasks = slice.Map(t.Tasks, func(task *Task) *Task {
		if task.Id == id {
			task.Status = status
		}
		return task
	})
}

func (t *TaskContainer) List(status ...Status) {
	if len(t.Tasks) == 0 {
		return
	}
	fmt.Printf("TaskContainer\n")
	fmt.Printf("%-4s\t%-24s\t%-8s\t%-10s\t%-10s\n", "ID", "Description", "Status", "CreatedAt", "UpdatedAt")
	tasks := slice.Filter(t.Tasks, func(task *Task) bool {
		if len(status) == 0 {
			return true
		}
		for _, s := range status {
			if s == task.Status {
				return true
			}
		}
		return false
	})
	slice.Range(tasks, func(_ int, task *Task) {
		fmt.Printf("%s\n", task)
	})
}

func (t *TaskContainer) Sync() {
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	write, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
		return
	}
	if write != len(bytes) {
		log.Fatal("failed to write all bytes")
	}
}

func main() {
	tasks, err := LoadFromFile(FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer tasks.Sync()
	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: %s add <description>\n", path.Base(os.Args[0]))
			return
		}
		tasks.Add(os.Args[2])
	case "list":
		if len(tasks.Tasks) == 0 {
			return
		}
		if len(os.Args) < 2 {
			fmt.Printf("Usage: %s list [status...]\n", path.Base(os.Args[0]))
			fmt.Printf("\tstatus\tdone,todo,in_progress\n")
			return
		}
		if len(os.Args) == 2 {
			tasks.List()
		}
		if len(os.Args) >= 3 {
			status := make([]Status, 0, len(os.Args[2:]))
			for _, s := range os.Args[2:] {
				s := Status(s)
				if s != Done && s != Todo && s != InProgress {
					fmt.Printf("%s is a invalid status\n", s)
					fmt.Printf("Usage: %s list [status...]\n", path.Base(os.Args[0]))
					fmt.Printf("\tstatus\tdone,todo,in_progress\n")
					return
				}
				status = append(status, s)
			}
			tasks.List(status...)
		}
	case "update":
		if len(tasks.Tasks) == 0 {
			return
		}
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s update <id> <description>\n", path.Base(os.Args[0]))
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid id %s", os.Args[2])
		}
		err = tasks.Update(id, os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	case "mark":
		if len(tasks.Tasks) == 0 {
			return
		}
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s mark <id> <status>\n", path.Base(os.Args[0]))
			fmt.Printf("\tstatus\tdone,todo,in_progress\n")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid id %s", os.Args[2])
		}
		s := Status(os.Args[3])
		if s != Done && s != Todo && s != InProgress {
			fmt.Printf("%s is a invalid status\n", s)
			fmt.Printf("Usage: %s list [status...]\n", path.Base(os.Args[0]))
			fmt.Printf("\tstatus\tdone,todo,in_progress\n")
			return
		}
		tasks.Mark(id, s)
	case "delete":
		if len(tasks.Tasks) == 0 {
			return
		}
		if len(os.Args) < 2 {
			fmt.Printf("Usage: %s delete <id>\n", path.Base(os.Args[0]))
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid id %s", os.Args[2])
		}
		err = tasks.Delete(id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
