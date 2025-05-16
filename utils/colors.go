package utils

// ANSI color codes for terminal output
const (
        // Reset all styles
        Reset = "\033[0m"
        
        // Regular colors
        Black  = "\033[30m"
        Red    = "\033[31m"
        Green  = "\033[32m"
        Yellow = "\033[33m"
        Blue   = "\033[34m"
        Purple = "\033[35m"
        Cyan   = "\033[36m"
        White  = "\033[37m"
        
        // Bold colors
        BoldBlack  = "\033[1;30m"
        BoldRed    = "\033[1;31m"
        BoldGreen  = "\033[1;32m"
        BoldYellow = "\033[1;33m"
        BoldBlue   = "\033[1;34m"
        BoldPurple = "\033[1;35m"
        BoldCyan   = "\033[1;36m"
        BoldWhite  = "\033[1;37m"
        
        // Background colors
        BgBlack  = "\033[40m"
        BgRed    = "\033[41m"
        BgGreen  = "\033[42m"
        BgYellow = "\033[43m"
        BgBlue   = "\033[44m"
        BgPurple = "\033[45m"
        BgCyan   = "\033[46m"
        BgWhite  = "\033[47m"
)

// ProgressState represents the state of a progress indicator
type ProgressState int

const (
        // Normal is the default state
        Normal ProgressState = iota
        // Success indicates a successful operation
        Success
        // Warning indicates a warning state
        Warning
        // Error indicates an error state
        Error
        // Info indicates an informational state
        Info
)

// State constants for backwards compatibility
const (
        StateNormal  = Normal
        StateSuccess = Success
        StateWarning = Warning
        StateError   = Error
        StateInfo    = Info
)

// GetColorForState returns the ANSI color code for a given progress state
func GetColorForState(state ProgressState) string {
        switch state {
        case Success:
                return Green
        case Warning:
                return Yellow
        case Error:
                return Red
        case Info:
                return Cyan
        default:
                return Reset
        }
}

// ColorizeText wraps text with the specified color and reset codes
func ColorizeText(text string, color string) string {
        return color + text + Reset
}