package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/st3iny/boards/internal/cli"
    "github.com/st3iny/boards/internal/task"
)

func main() {
    tasks, err := task.LoadStorage()
    if err != nil {
        fmt.Println("Error while trying to load storage.json")
        panic(err)
    }

    args := os.Args[1:]
    var cmdArgs []string
    if len(args) > 1 {
        cmdArgs = args[1:]
    }

    var cliError error
    if len(args) == 0 {
        cli.BoardView(tasks)
    } else if args[0] == "t" || args[0] == "task" {
        cliError = cli.CreateTask(tasks, cmdArgs)
    } else if args[0] == "n" || args[0] == "note" {
        cliError = cli.CreateNote(tasks, cmdArgs)
    } else if args[0] == "b" || args[0] == "boards" {
        cliError = cli.EditBoards(tasks, cmdArgs)
    } else if args[0] == "e" || args[0] == "edit" {
        cliError = cli.EditDescription(tasks, cmdArgs)
    } else if args[0] == "c" || args[0] == "complete" {
        cliError = cli.Complete(tasks, cmdArgs)
    } else if args[0] == "u" || args[0] == "uncomplete" {
        cliError = cli.Uncomplete(tasks, cmdArgs)
    } else if args[0] == "d" || args[0] == "delete" {
        cliError = cli.Delete(tasks, cmdArgs)
    } else if args[0] == "clear" {
        cliError = cli.Clear(tasks)
    } else if args[0] == "-h" || args[0] == "--help" {
        usage()
        os.Exit(0)
    } else {
        fmt.Println("Invalid command (try --help)")
        os.Exit(1)
    }

    if cliError != nil {
        fmt.Println("Error while trying to save storage.json")
        panic(cliError)
    }
}

func usage() {
    help := []string{
        "boards",
        "boards task|t [@BOARD ...] DESCRIPTION [DESCRIPTION ...]",
        "boards note|n [@BOARD ...] DESCRIPTION [DESCRIPTION ...]",
        "boards delete|d ID [ID ...]",
        "boards complete|c ID [ID ...]",
        "boards uncomplete|u ID [ID ...]",
        "boards edit|e ID DESCRIPTION [DESCRIPTION ...]",
        "boards boards|b ID [ID ...] [@BOARD ...]",
        "boards clear",
        "boards --help|-h",
        "",
        "If boards is run without any argument all tasks will be printed grouped by their boards.",
        "Clear will remove all complete tasks.",
    }
    fmt.Println(strings.Join(help, "\n"))
}
