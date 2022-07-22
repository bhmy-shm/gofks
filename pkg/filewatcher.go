package pkg

import (
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

type watcher struct {
	f *File

	fw   *fsnotify.Watcher
	exit chan bool
}

func newWatcher(f *File) (Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(f.path)

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*File, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, errorx.WatcherFileStop
	default:
	}

	select {
	case event, _ := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			// check existence of File, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}

		//重新读取加载文件，并返回文件
		return w.f.Reload(), nil

	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, errorx.WatcherFileStop
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}

func ReadWatcher(f *File) {
	w, _ := f.Watch()
	newf, err := w.Next()

	if err != nil {
		log.Fatalln("watch faile =", err)
	}
	f = newf
	if f.confMap != nil {
		ReadWatcher(f)
	} else {
		log.Fatalln("监听后的内容出现异常")
	}
}
