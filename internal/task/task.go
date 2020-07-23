package task

import (
    "fmt"
    "strconv"

    "github.com/st3iny/boards/internal/style"
)

const DefaultBoard = "My Board"

type Task struct {
    Id int `json:"id"`
    Description string `json:"description"`
    Boards []string `json:"boards"`
    Note bool `json:"note"`
    Done bool `json:"complete"`
    Urgent bool `json:"urgent"`
}

func (task Task) Print(indent int) {
    id := strconv.FormatInt(int64(task.Id), 10)
    for len(id) < indent {
        id = " " + id
    }
    id = style.Muted + id + "." + style.Reset

    var tick string
    if task.Note {
        tick = style.Note + "●" + style.Reset
    } else if task.Done {
        tick = style.Done + "✔" + style.Reset
    } else {
        tick = style.Pending + "☐" + style.Reset
    }

    description := task.Description
    if !task.Note && task.Done {
        description = style.Muted + task.Description + style.Reset
    } else if task.Urgent {
        description = style.Urgent + task.Description + style.Reset
    }

    fmt.Printf("  %s %s  %s\n", id, tick, description)
}
