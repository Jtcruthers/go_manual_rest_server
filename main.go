package main 

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strconv"
  "strings"
  "time"

  "jtcruthers/manual_rest_server/internal/taskstore"
)

type taskServer struct {
    store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
    store := taskstore.New()
    return &taskServer{store: store}
}

func renderJson(w http.ResponseWriter, v interface{}) {
  js, err := json.Marshal(v)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func (ts *taskServer) taskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTasksHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE, or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
      path := strings.Trim(req.URL.Path, "/")
      pathParts := strings.Split(path, "/")
      if len(pathParts) < 2 {
        http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
        return
      }

      id, err := strconv.Atoi(pathParts[1])
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

      if req.Method == http.MethodDelete {
        ts.deleteTaskHandler(w, req, int(id))
      } else if req.Method == http.MethodGet {
        ts.getTaskHandler(w, req, int(id))
      } else {
        http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
        return
      }
    }
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {  	
  log.Printf("handling get all tasks at %s\n", req.URL.Path)

  allTasks := ts.store.GetAllTasks()

  renderJson(w, allTasks)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request, id int) {  	
  log.Printf("handling get all tasks at %s\n", req.URL.Path)

  task, err := ts.store.GetTask(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  renderJson(w, task)
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
  log.Printf("handling delete all tasks at %s\n", req.URL.Path)

  err := ts.store.DeleteAllTasks()
  if err !=  nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
  log.Printf("handling task delete at %s\n", req.URL.Path)

  err := ts.store.DeleteTask(id)
  if err !=  nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

type CreateTaskRequest struct {
  Text string
  Tags []string
  Due time.Time
}

type IdResponse struct {
  Id int
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
  log.Printf("handling create task at %s\n", req.URL.Path)
  
  var createTask CreateTaskRequest
  err := json.NewDecoder(req.Body).Decode(&createTask)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  log.Println(time.Now())
  id := ts.store.CreateTask(createTask.Text, createTask.Tags, createTask.Due)

  res := IdResponse{Id: id}
  renderJson(w, res)
}

func (ts *taskServer) getTasksByTagHandler(w http.ResponseWriter, req *http.Request, tag string) {
  log.Printf("handling get tasks by tag at %s\n", req.URL.Path)

  tasks := ts.store.GetTasksByTag(tag)
  renderJson(w, tasks)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/tag/" {
      http.Error(w, "expect /tag/<tagname>, got /tag/", http.StatusBadRequest)
      return
	} else {
      path := strings.Trim(req.URL.Path, "/")
      pathParts := strings.Split(path, "/")
      if len(pathParts) < 2 {
        http.Error(w, "expect only /tag/<tagname> in task handler", http.StatusBadRequest)
        return
      }

      tag := pathParts[1]

      if req.Method == http.MethodGet {
        ts.getTasksByTagHandler(w, req, tag)
      } else {
        http.Error(w, fmt.Sprintf("expect method GET at /tag/<tagname>, got %v", req.Method), http.StatusMethodNotAllowed)
        return
      }
    }
}

func main() {
    server := NewTaskServer()
    mux := http.NewServeMux()
    mux.HandleFunc("/task/", server.taskHandler)
    mux.HandleFunc("/tag/", server.tagHandler)

    log.Fatal(http.ListenAndServe("localhost:4000", mux))
}

