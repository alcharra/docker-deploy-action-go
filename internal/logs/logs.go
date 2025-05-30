package logs

import (
	"fmt"
	"os"
)

var IsVerbose bool

func Step(title string) {
	fmt.Println()
	fmt.Println(title)
}

func Stepf(format string, args ...interface{}) {
	fmt.Println()
	fmt.Printf(format+"\n", args...)
}

func Substep(msg string) {
	fmt.Printf("   %s\n", msg)
}

func Substepf(format string, args ...interface{}) {
	fmt.Printf("   "+format+"\n", args...)
}

func Warn(msg string) {
	fmt.Printf("   \U000026A0\U0000FE0F  %s\n", msg)
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf("   \U000026A0\U0000FE0F  "+format+"\n", args...)
}

func Error(msg string) {
	fmt.Printf("   \U0000274C %s\n", msg)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf("   \U0000274C "+format+"\n", args...)
}

func Info(msg string) {
	fmt.Printf("   \U00002139\U0000FE0F  %s\n", msg)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf("   \U00002139\U0000FE0F  "+format+"\n", args...)
}

func Success(msg string) {
	fmt.Printf("   \U00002705 %s\n", msg)
}

func Successf(format string, args ...interface{}) {
	fmt.Printf("   \U00002705 "+format+"\n", args...)
}

func Verbose(msg string) {
	if IsVerbose {
		fmt.Printf("   \U0001F50D %s\n", msg)
	}
}

func Verbosef(format string, args ...interface{}) {
	if IsVerbose {
		fmt.Printf("   \U0001F50D "+format+"\n", args...)
	}
}

func VerboseCommand(cmd string) {
	if IsVerbose {
		fmt.Printf("      \U000027A5 Running command: %s\n", cmd)
	}
}

func VerboseCommandf(format string, args ...interface{}) {
	if IsVerbose {
		fmt.Printf("      \U000027A5 Running command: "+format+"\n", args...)
	}
}

func Fatal(msg string) {
	fmt.Println()
	fmt.Printf("\U0000274C %s\n", msg)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Println()
	msg := fmt.Sprintf("\U0000274C "+format+"\n", args...)
	fmt.Print(msg)
	os.Exit(1)
}

func Break() {
	fmt.Println()
}
