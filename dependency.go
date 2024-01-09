package qpm

import (
	"context"
	"slices"
	"sync"
)

// dependencies stratumの依存しているstratumを再帰的に取得する
func dependencies(c Config, a Action, os OS, stratumName string, knownDeps map[string][]string) error {
	sf, err := readStratumFile(c.AquiferPath, stratumName)
	if err != nil {
		return err
	}

	jobs := sf[a.String()]

	var deps []string
	for _, v := range jobs {
		if slices.Contains(v.OS, os.String()) {
			deps = v.Dependency
			break
		}
	}

	knownDeps[stratumName] = deps

	for _, v := range deps {
		if err := dependencies(c, a, os, v, knownDeps); err != nil {
			return err
		}
	}

	return nil
}

type task struct {
	pkg  string
	deps []string
}

type multiTaskExec struct {
	wg       *sync.WaitGroup
	started  bool
	tasks    []task
	packages *sync.Map
}

func (m *multiTaskExec) add(pkg string, deps []string) {
	if m.started {
		panic("already started")
	}

	m.wg.Add(1)

	m.tasks = append(m.tasks, task{pkg: pkg, deps: deps})

	if m.packages == nil {
		m.packages = &sync.Map{}
	}
	if _, ok := m.packages.Load(pkg); !ok {
		m.packages.Store(pkg, make(chan struct{}))
	}
}

func (*multiTaskExec) waitAllChans(chans []chan struct{}) {
	for _, c := range chans {
		for {
			if _, ok := <-c; !ok {
				break
			}
		}
	}
}

func (m *multiTaskExec) execTask(ctx context.Context, t task, f func(string)) {
	depChans := make([]chan struct{}, 0)
	for _, dep := range t.deps {
		c, ok := m.packages.Load(dep)
		if ok {
			depChans = append(depChans, c.(chan struct{}))
		}
	}

	m.waitAllChans(depChans)

	f(t.pkg)

	c, ok := m.packages.Load(t.pkg)
	if ok {
		close(c.(chan struct{}))
		m.packages.Delete(t.pkg)
	}
}

func (m *multiTaskExec) wait(f func(string)) {
	if m.started {
		panic("already started")
	}

	m.started = true

	for _, t := range m.tasks {
		go func(t task) {
			defer m.wg.Done()
			m.execTask(context.TODO(), t, f)
		}(t)
	}

	m.wg.Wait()
}

func newMultiTaskExec() *multiTaskExec {
	var wg sync.WaitGroup
	return &multiTaskExec{wg: &wg}
}
