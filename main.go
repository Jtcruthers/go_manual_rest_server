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
  js, err := json.Marshal(allTasks)
  if err !=  nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request, id int) {  	
  log.Printf("handling get all tasks at %s\n", req.URL.Path)

  task, err := ts.store.GetTask(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  js, err := json.Marshal(task)
  if err !=  nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
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

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
  log.Printf("handling create task at %s\n", req.URL.Path)
  
  text := "Hey" // TODO - fix
  tags := []string{"tag1", "tag2"}
  due := time.Now()
  ts.store.CreateTask(text, tags, due)
}

func main() {
    server := NewTaskServer()
    mux := http.NewServeMux()
    mux.HandleFunc("/task/", server.taskHandler)
    //mux.HandleFunc("/tag/", server.tagHandler)
    //mux.HandleFunc("/due/", server.dueHandler)

    log.Fatal(http.ListenAndServe("localhost:4000", mux))
}

