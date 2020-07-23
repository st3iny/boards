package style

const (
    Reset   = "\033[0m"
    Board   = "\033[4m"        // underlined
    Urgent  = "\033[38;5;3m"   // yellow
    Muted   = "\033[38;5;7m"   // gray
    Pending = "\033[38;5;5m"   // purple
    Done    = "\033[38;5;2m"   // green
    Note    = "\033[38;5;4m"   // blue
    Add     = "\033[4;38;5;2m" // underlined, green
    Remove  = "\033[4;38;5;1m" // underlined, red
)
