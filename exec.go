package gitls

import (
	"regexp"
	"strings"
)

// cutMessage get remote, pull, push branches
/*
 * return
 *	- First  -> remote branches
 *	- Second -> pull branches
 *	- Third  -> push branches
 */
func cutMessage(message string) map[string][]string {
	remoteBranch := "  Remote branch:"
	remoteBranches := "  Remote branches:"
	gitPull := "  Local branch configured for 'git pull':"
	gitPulls := "  Local branches configured for 'git pull':"
	gitPush := "  Local ref configured for 'git push':"
	gitPushes := "  Local refs configured for 'git push':"

	cutMark := []string{
		remoteBranch, remoteBranches,
		gitPull, gitPulls,
		gitPush, gitPushes}

	// split message to 3 substring index
	mLines := strings.Split(message, "\n")
	cutNum := 0
	cutIndex := []int{-1, -1, -1}
	for k, v := range mLines {
		for i := cutNum; i < len(cutMark); i++ {
			if v == cutMark[i] {
				if i == 0 || i == 1 {
					cutIndex[0] = k
				}
				if i == 2 || i == 3 {
					cutIndex[1] = k
				}
				if i == 4 || i == 5 {
					cutIndex[2] = k
				}
				cutNum += 2
				break
			}
		}
	}

	result := make(map[string][]string, 0)
	result["push"] = make([]string, 0)
	result["pull"] = make([]string, 0)
	result["remote"] = make([]string, 0)
	// last line of message maybe blank -> do nothing of any blank line
	// may reverse order of result
	for i := len(mLines) - 1; i >= 0; i-- {
		if mLines[i] == "" {
			continue
		}

		m := strings.TrimSpace(mLines[i])

		if i == cutIndex[2] {
			cutIndex[2] = -1
			continue
		}
		if i > cutIndex[2] && cutIndex[2] != -1 {
			result["push"] = append(result["push"], m)
			continue
		}
		if i == cutIndex[1] {
			cutIndex[1] = -1
			continue
		}
		if i > cutIndex[1] && cutIndex[1] != -1 {
			result["pull"] = append(result["pull"], m)
			continue
		}
		if i == cutIndex[0] {
			cutIndex[0] = -1
			continue
		}
		if i > cutIndex[0] && cutIndex[0] != -1 {
			result["remote"] = append(result["remote"], m)
			continue
		}
	}

	return result
}

// GetPushStatus get branch status for push
func GetPushStatus(message string) map[string]string {
	words := strings.Fields(message)
	branch := words[0]
	rgx := regexp.MustCompile(`\((.*?)\)`)
	result := make(map[string]string, 0)
	result[branch] = strings.Trim(rgx.FindString(message), "()")

	return result
}
