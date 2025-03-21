package colorfmt

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
	Magenta = "\033[35m"
	Reset  = "\033[0m"
)

// PrintColor prints formatted text in the specified color
func PrintColor(color, format string, a ...any) (int, error) {
	return fmt.Printf("%s%s%s", color, fmt.Sprintf(format, a...), Reset)
}

// Convenience functions for common colors
func PrintRed(format string, a ...any) (int, error) {
	return PrintColor(Red, format, a...)
}

func PrintGreen(format string, a ...any) (int, error) {
	return PrintColor(Green, format, a...)
}

func PrintBlue(format string, a ...any) (int, error) {
	return PrintColor(Blue, format, a...)
}

func PrintYellow(format string, a ...any) (int, error) {
	return PrintColor(Yellow, format, a...)
}


func Sprintf(color, format string, a ...any) string {
	return fmt.Sprintf("%s%s%s", color, fmt.Sprintf(format, a...), Reset)
}
