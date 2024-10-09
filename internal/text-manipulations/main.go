package textmanipulations

import "strings"

func Remove_prefixes(content []string, prefixes []string) []string {
	var result []string

	for i := 0; i < len(content); i++ {
		line := content[i]
		found := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				found = true
				i = remove_following_block(get_identation(line), i, content)
				break
			}
		}
		if !found {
			result = append(result, line)
		}
	}

	return result
}

func Remove_lines(content []string, prefixes []string) []string {
	var result []string

	for i := 0; i < len(content); i++ {
		line := content[i]
		found := false
		for _, prefix := range prefixes {
			if line == prefix {
				found = true
				i = remove_following_block(get_identation(line), i, content)
				break
			}
		}
		if !found {
			result = append(result, line)
		}
	}

	return result
}

func remove_following_block(identation int, i int, content []string) int {
	block_identation := strings.Repeat(" ", identation+1)
	for i < len(content)-1 && strings.HasPrefix(content[i+1], block_identation) {
		i++
	}

	return i
}

func get_identation(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
