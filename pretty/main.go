package pretty

import "fmt"

func Up(n int) {
	fmt.Printf("\033[%dA", n)
}

func Down(n int) {
	fmt.Printf("\033[%dB", n)
}

func Right(n int) {
	fmt.Printf("\033[%dC", n)
}

func Left(n int) {
	fmt.Printf("\033[%dD", n)
}

func Erase() {
	fmt.Printf("\033[2K\033[G")
}

func Bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}

func UnderLine(str string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", str)
}

func Black(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Red(str string) string {
	return fmt.Sprintf("\033[31m%s\033[39m", str)
}

func Green(str string) string {
	return fmt.Sprintf("\033[32m%s\033[39m", str)
}

func Yellow(str string) string {
	return fmt.Sprintf("\033[33m%s\033[39m", str)
}

func Blue(str string) string {
	return fmt.Sprintf("\033[34m%s\033[39m", str)
}

func Magenta(str string) string {
	return fmt.Sprintf("\033[35m%s\033[39m", str)
}

func Cyan(str string) string {
	return fmt.Sprintf("\033[36m%s\033[39m", str)
}

func White(str string) string {
	return fmt.Sprintf("\033[37m%s\033[39m", str)
}
