package worktracker

import (
	"fmt"
	"io"
	"strings"
)

func printUsage(p string) error {
	return fmt.Errorf("USAGE: %s start|stop|get [type]", p)
}

func HandleCommand(w io.Writer, a []string, s *WorkService) error {
	if len(a) < 2 {
		return printUsage(a[0])
	}

	args := a[1:]
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
				PrintWorks(w, works)
			}
			return err
		}
		works, err := s.GetWork(ID)
		if works != nil {
			PrintWorks(w, works)
		}
		return err
	default:
		return printUsage(a[0])
	}
}
