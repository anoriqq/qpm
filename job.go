package qpm

type step struct {
	name string
	run  string
}

type job struct {
	dependency     []string
	availableShell map[string]struct{}
	step           []step
}

func (j job) shell(prioritizedShells []string) (string, bool) {
	if len(prioritizedShells) == 0 {
		return "", false
	}
	if len(j.availableShell) == 0 {
		return "", false
	}

	for _, v := range prioritizedShells {
		if _, ok := j.availableShell[v]; ok {
			return v, true
		}
	}

	return "", false
}
