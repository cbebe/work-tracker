package worktracker

import (
	"fmt"
	"os"
	"strings"
)

func HandleCommand(args []string, s *WorkService) error {
	switch strings.ToLower(args[0]) {
	case "start":
		if len(args) >= 2 {
			return s.StartLog(args[1], ID)
		}
		return s.StartWork()
	case "stop":
		if len(args) >= 2 {
			return s.StopLog(args[1], ID)
		}
		return s.StopWork()
	case "list":
		fallthrough
	case "get":
		if len(args) >= 2 {
			works, err := s.GetWorkType(args[1], ID)
			if works != nil {
				PrintWorks(os.Stdout, works)
			}
			return err
		}
		works, err := s.GetWork(ID)
		if works != nil {
			PrintWorks(os.Stdout, works)
		}
		return err
	default:
		PrintUsage()
		return nil
	}
}

func PrintUsage() {
	fmt.Println("USAGE:", os.Args[0], "start|stop|get [type]")
}
