package chatgpt

import (
	"fmt"
)

func TransDecorate(content, language string) (string, error) {
	switch language {
	case English:
		return fmt.Sprintf("Translate the following into English:\n%q", content), nil
	case Chinese:
		return fmt.Sprintf("Translate the following into Chinese:\n%q", content), nil
	default:
		return "", fmt.Errorf("unknown languge")
	}
}
