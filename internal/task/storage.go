package task

import (
    "encoding/json"
    "io/ioutil"
    "os"

    "github.com/adrg/xdg"
)

func getStoragePath() (string, error) {
    storagePath, err := xdg.ConfigFile("boards/storage.json")
    if err != nil {
        return "", err
    }

    return storagePath, nil
}

func LoadStorage() ([]Task, error) {
    storagePath, err := getStoragePath()
    if err != nil {
        return nil, err
    }

    _, err = os.Stat(storagePath)
    if os.IsNotExist(err) {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    tasksBlob, err := ioutil.ReadFile(storagePath)
    if err != nil {
        return nil, err
    }

    var tasks []Task
    err = json.Unmarshal(tasksBlob, &tasks)
    if err != nil {
        return nil, err
    }

    return tasks, nil
}

type Storage []Task

func (tasks Storage) Save() error {
    tasksBlob, err := json.Marshal(tasks)
    if err != nil {
        return err
    }

    storagePath, err := getStoragePath()
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(storagePath, tasksBlob, 0600)
    if err != nil {
        return err
    }

    return nil
}

func (tasks Storage) GetHighestId() int {
    highestId := 0
    for _, item := range tasks {
        if item.Id > highestId {
            highestId = item.Id
        }
    }

    return highestId
}

func (tasks Storage) GetNextId() int {
    freeId := 1
    for {
        used := false

        for _, item := range tasks {
            if freeId == item.Id {
                used = true
                break
            }
        }

        if !used {
            break
        }
        freeId++
    }

    return freeId
}

func (tasks Storage) ContainsTask(id int) bool {
    for _, item := range tasks {
        if item.Id == id {
            return true
        }
    }

    return false
}

func (tasks Storage) GetTask(id int) *Task {
    for index, item := range tasks {
        if item.Id == id {
            return &tasks[index]
        }
    }

    return nil
}

func (tasks *Storage) DeleteTask(id int) bool {
    index := -1
    for i, item := range *tasks {
        if item.Id == id {
            index = i
            break
        }
    }

    if index != -1 {
        *tasks = append((*tasks)[:index], (*tasks)[index+1:]...)
        return true
    }

    return false
}
