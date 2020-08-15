package pretty

import "fmt"

func Up(n int) {
	fmt.Printf("\033[%dA", n)
}

func Down(n int) {
	fmt.Printf("\033[%dA", n)
}

func Left(n int) {
	fmt.Printf("\033[%dA", n)
}

func Right(n int) {
	fmt.Printf("\033[%dA", n)
}

func Erase() {
	fmt.Printf("\033[2K\033[G")
}

func Black(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Red(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Green(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Yellow(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Blue(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Magenta(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func Cyan(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}

func White(str string) string {
	return fmt.Sprintf("\033[30m%s\033[39m", str)
}
