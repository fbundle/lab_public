package feature_toggle

type Loader interface {
	Load(key interface{}) bool
}

type Toggle interface {
	Exec()
	If(key interface{}, task func(), alt func()) Toggle
}

func Set(loader Loader) {
	defaultToggle.loader = loader
}

func If(key interface{}, task func(), alt func()) Toggle {
	return defaultToggle.If(key, task, alt)
}

var defaultToggle = &toggle{
	loader: nil,
	parent: nil,
	task:   nil,
}

type toggle struct {
	loader Loader
	parent *toggle
	task   func()
}

func (t *toggle) Exec() {
	if t.parent != nil {
		t.parent.Exec()
	}
	if t.task != nil {
		t.task()
	}
}

func (t *toggle) If(key interface{}, task func(), alt func()) Toggle {
	chosen := alt
	if t.loader != nil && t.loader.Load(key) {
		chosen = task
	}
	return &toggle{
		loader: t.loader,
		parent: t,
		task:   chosen,
	}
}
