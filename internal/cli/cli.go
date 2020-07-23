package cli

import (
    "fmt"
    "sort"
    "strings"
    "strconv"

    "github.com/st3iny/boards/internal/style"
    "github.com/st3iny/boards/internal/task"

    "github.com/jwangsadinata/go-multimap/slicemultimap"
)

func BoardView(tasks task.Storage) {
    done := 0
    pending := 0
    notes := 0
    boardMap := slicemultimap.New()
    for _, item := range tasks {
        if item.Note {
            notes++
        } else if item.Done {
            done++
        } else {
            pending++
        }

        if len(item.Boards) == 0 {
            boardMap.Put(task.DefaultBoard, item)
        } else {
            for _, board := range item.Boards {
                boardMap.Put("@" + board, item)
            }
        }
    }

    keys := boardMap.KeySet()
    boards := make([]string, len(keys))
    for index, board := range keys {
        boards[index] = board.(string)
    }
    sort.Strings(boards)

    indent := len(strconv.FormatInt(int64(tasks.GetHighestId()), 10))
    for _, board := range boards {
        fmt.Printf("\n%s%s%s\n", style.Board, board, style.Reset)
        tasks, _ := boardMap.Get(board)
        for _, item := range tasks {
            item.(task.Task).Print(indent)
        }
    }

    fmt.Printf(
        "\n%s%d%s done, %s%d%s pending, %s%d%s notes%s\n",
        style.Done, done, style.Muted,
        style.Pending, pending, style.Muted,
        style.Note, notes, style.Muted,
        style.Reset,
    )
}

func createItem(tasks task.Storage, note bool, args []string) (int, error) {
    if len(args) == 0 {
        return 0, nil
    }

    var description []string
    var boards []string
    for _, arg := range args {
        if strings.HasPrefix(arg, "@") {
            boards = append(boards, strings.TrimPrefix(arg, "@"))
        } else {
            description = append(description, arg)
        }
    }

    id := tasks.GetNextId()
    tasks = append(tasks, task.Task{
        Id: id,
        Description: strings.Join(description, " "),
        Boards: boards,
        Note: note,
        Done: false,
    })

    return id, tasks.Save()
}

func CreateTask(tasks task.Storage, args []string) error {
    id, err := createItem(tasks, false, args)
    if id != 0 {
        fmt.Printf("Created task %s%d%s\n", style.Muted, id, style.Reset)
    }

    return err
}

func CreateNote(tasks task.Storage, args []string) error {
    id, err := createItem(tasks, true, args)
    if id != 0 {
        fmt.Printf("Created note %s%d%s\n", style.Muted, id, style.Reset)
    }

    return err
}

func EditBoards(tasks task.Storage, args []string) error {
    if len(args) == 0 {
        return nil
    }

    deltaMode := false
    var ids []int
    var boards, removeBoards []string
    for _, arg := range args {
        if !deltaMode && (strings.HasPrefix(arg, "+") || strings.HasPrefix(arg, "-")) {
            deltaMode = true
        }

        if strings.HasPrefix(arg, "@") || strings.HasPrefix(arg, "+@") {
            board := strings.TrimPrefix(arg, "+")
            boards = append(boards, strings.TrimPrefix(board, "@"))
        } else if strings.HasPrefix(arg, "-@") {
            removeBoards = append(removeBoards, strings.TrimPrefix(arg, "-@"))
        } else {
            id, err := strconv.Atoi(arg)
            if err == nil && tasks.ContainsTask(id) {
                ids = append(ids, id)
            }
        }
    }

    if len(ids) == 0 {
        return nil
    }

    for _, id := range ids {
        task := tasks.GetTask(id)
        if deltaMode {
            task.Boards = append(task.Boards, boards...)
            for _, remove := range removeBoards {
                for index, board := range task.Boards {
                    if board == remove {
                        task.Boards = append(task.Boards[:index], task.Boards[index + 1:]...)
                        break
                    }
                }
            }
        } else {
            task.Boards = boards
        }
    }

    idsDisplay := strings.Join(surround(stringify(ids), style.Muted, style.Reset), ", ")
    var boardsDisplay string
    if deltaMode {
        greenBoards := surround(boards, style.Add + "+@", style.Reset)
        redBoards := surround(removeBoards, style.Remove + "-@", style.Reset)
        boardsDisplay = strings.Join(append(greenBoards, redBoards...), ", ")
    } else {
        boardsDisplay = strings.Join(surround(boards, style.Board + "@", style.Reset), ", ")
    }
    fmt.Printf("Changed boards of items [ %s ] to [ %s ]\n", idsDisplay, boardsDisplay)
    return tasks.Save()
}

func EditDescription(tasks task.Storage, args []string) error {
    if len(args) < 2 {
        return nil
    }

    id, err := strconv.Atoi(args[0])
    if err != nil || !tasks.ContainsTask(id) {
        return nil
    }

    description := strings.Join(args[1:], " ")
    tasks.GetTask(id).Description = description
    fmt.Printf(
        "Changed description of item %s%d%s to \"%s\"\n",
        style.Muted, id, style.Reset, description,
    )
    return tasks.Save()
}

func setDone(tasks task.Storage, done bool, args []string) ([]string, error) {
    if len(args) == 0 {
        return nil, nil
    }

    var ids []int
    for _, arg := range args {
        id, err := strconv.Atoi(arg)
        if err != nil || !tasks.ContainsTask(id) {
            continue
        }

        task := tasks.GetTask(id)
        if !task.Note && task.Done != done {
            task.Done = done
            ids = append(ids, id)
        }
    }

    if len(ids) == 0 {
        return nil, nil
    }

    return surround(stringify(ids), style.Muted, style.Reset), tasks.Save()
}

func Complete(tasks task.Storage, args []string) error {
    ids, err := setDone(tasks, true, args)
    if err == nil && len(ids) > 0 {
        fmt.Printf("Completed tasks [ %s ]\n", strings.Join(ids, ", "))
    }

    return err
}

func Uncomplete(tasks task.Storage, args []string) error {
    ids, err := setDone(tasks, false, args)
    if err == nil && len(ids) > 0 {
        fmt.Printf("Uncompleted tasks [ %s ]\n", strings.Join(ids, ", "))
    }

    return err
}

func ToggleUrgency(tasks task.Storage, args []string) error {
    if len(args) == 0 {
        return nil
    }

    var ids []int
    for _, arg := range args {
        id, err := strconv.Atoi(arg)
        if err != nil || !tasks.ContainsTask(id) {
            continue
        }

        task := tasks.GetTask(id)
        task.Urgent = !task.Urgent
        ids = append(ids, id)
    }

    if len(ids) == 0 {
        return nil
    }

    idsDisplay := strings.Join(surround(stringify(ids), style.Muted, style.Reset), ", ")
    fmt.Printf("Toggled urgency of items [ %s ]\n", idsDisplay)
    return tasks.Save()
}

func Delete(tasks task.Storage, args []string) error {
    if len(args) == 0 {
        return nil
    }

    var ids []int
    for _, arg := range args {
        id, err := strconv.Atoi(arg)
        if err != nil || !tasks.ContainsTask(id) {
            continue
        }

        ids = append(ids, id)
    }

    if len(ids) == 0 {
        return nil
    }

    for _, id := range ids {
        tasks.DeleteTask(id)
    }

    idsDisplay := strings.Join(surround(stringify(ids), style.Muted, style.Reset), ", ")
    fmt.Printf("Deleted items [ %s ]\n", idsDisplay)
    return tasks.Save()
}

func Clear(tasks task.Storage) error {
    var ids []int
    for _, item := range tasks {
        if !item.Note && item.Done {
            ids = append(ids, item.Id)
        }
    }

    if len(ids) == 0 {
        return nil
    }

    for _, id := range ids {
        tasks.DeleteTask(id)
    }

    idsDisplay := strings.Join(surround(stringify(ids), style.Muted, style.Reset), ", ")
    fmt.Printf("Deleted tasks [ %s ]\n", idsDisplay)
    return tasks.Save()
}
