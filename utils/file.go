package sgc7utils

import "io/ioutil"

// IsSameFile - filea == fileb
func IsSameFile(fna string, fnb string) bool {
	da, err := ioutil.ReadFile(fna)
	if err != nil {
		return false
	}

	db, err := ioutil.ReadFile(fnb)
	if err != nil {
		return false
	}

	if len(da) != len(db) {
		return false
	}

	for i, v := range da {
		if v != db[i] {
			return false
		}
	}

	return true
}
