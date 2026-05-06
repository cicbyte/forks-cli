package clone

import "fmt"

func validateCloneParams(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("必须指定仓库地址")
	}
	return nil
}
