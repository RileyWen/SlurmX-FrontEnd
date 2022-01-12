package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type SbatchArg struct {
	name string
	val  string
}

func ProcessSbatchArg(args []SbatchArg) {

}

func ProcessLine(line string, sh *[]string, args *[]SbatchArg) bool {
	re := regexp.MustCompile(`^#SBATCH`)
	if re.MatchString(line) {
		split := strings.Fields(line)
		if len(split) == 3 {
			*args = append(*args, SbatchArg{name: split[1], val: split[2]})
		} else if len(split) == 2 {
			*args = append(*args, SbatchArg{name: split[1]})
		} else {
			return false
		}
	} else {
		*sh = append(*sh, line)
	}

	return true
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	num := 0
	sh := make([]string, 0)
	args := make([]SbatchArg, 0)

	for scanner.Scan() {
		num++
		success := ProcessLine(scanner.Text(), &sh, &args)
		if !success {
			err = fmt.Errorf("grammer error at line %v", num)
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Invoking UID: %d\n\n", os.Getuid())

	fmt.Printf("Shell script:\n%s\n\n", strings.Join(sh, "\n"))

	fmt.Printf("Sbatch args:\n%v\n", args)
}
