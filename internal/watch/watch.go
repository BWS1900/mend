package watch

import (
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	w *fsnotify.Watcher
}

func New(path string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(path); err != nil {
		w.Close()
		return nil, err
	}
	return &Watcher{w: w}, nil
}

func (w *Watcher) Events() <-chan bool {
	out := make(chan bool, 1)
	go func() {
		defer close(out)
		for {
			select {
			case ev, ok := <-w.w.Events:
				if !ok {
					return
				}
				if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
					select {
					case out <- true:
					default:
					}
				}
			case _, ok := <-w.w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	return out
}

func (w *Watcher) Close() error {
	return w.w.Close()
}
