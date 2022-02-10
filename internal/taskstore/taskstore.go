package taskstore

import (
  "fmt"
  "time"
)

type Task struct {
  Id   int       `json:"id"`
  Text string    `json:"text"`
  Tags []string  `json:"tags"`
  Due  time.Time `json:"due"`
}

type TaskStore struct {
  tasks map[int]Task
  nextId int
}

func New() *TaskStore {
  ts := &TaskStore{}
  ts.tasks = make(map[int]Task)
  ts.nextId = 0

  return ts
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int  {
  task := Task{
    Id: ts.nextId,
    Text: text,
    Due: due,
    Tags: tags,
  }

  ts.tasks[ts.nextId] = task
  ts.nextId++

  return task.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
  task, ok := ts.tasks[id]

  if ok {
    return task, nil
  }

  return Task{}, fmt.Errorf("task with id=%d not found", id)
}

func (ts *TaskStore) DeleteTask(id int) error {
  _, ok := ts.tasks[id]

  if !ok {
    return fmt.Errorf("task with id=%d not found for deletion", id)
  }

  delete(ts.tasks, id)
  return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
  ts.tasks = make(map[int]Task)   
  ts.nextId = 0

  return nil
}

func (ts *TaskStore) GetAllTasks() []Task {
  allTasks := make([]Task, 0, len(ts.tasks))

  for _, task := range ts.tasks {
    allTasks = append(allTasks, task)
  }

  return allTasks
}

//func (ts *TaskStore) GetTasksByTag(tag string) []Task

//func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task

