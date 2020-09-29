package gitls

import (
	"os"
)

// IsFile check file exist
/* IsFile
 *  - require absolute path of file
 *  - file exist 			-> true, nil
 *  - file not exist	-> false, nil
 *  - error 					-> _, err
 */
func IsFile(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if !os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

// IsDir check if is a Dir
/* IsDir
 *  - require absolute path of Dir
 *  - dir exist 			-> true, nil
 *  - dir not exist		-> false, nil
 *  - error 					-> _, err
 */
func IsDir(dir string) (bool, error) {
	d, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if !os.IsNotExist(err) {
			return false, err
		}
	}
	return d.IsDir(), nil
}