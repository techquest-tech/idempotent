package idempotent

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/techquest-tech/gobatch"
)

type PersistedIdempotent struct {
	File     string `default:".idempotent.txt"`
	Batch    uint   `default:"100"`
	Interval string `default:"10s"`
	Retry    uint   `default:"3"`
	Service  IdempotentService
	ch       chan interface{}
}

func (f *PersistedIdempotent) Init() error {
	if f.Service == nil {
		return fmt.Errorf("parent service can't be empty, either InMemoryMap or Radix")
	}

	defaults.Set(f)

	file, err := os.OpenFile(f.File, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		msg := fmt.Sprintf("create or open file error, file = %s, err = %v", f.File, err)
		log.Error(msg)
		return err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	// initCached := map[interface{}]bool{}
	for sc.Scan() {
		// lines= append(lines, sc.Text())
		// initCached[sc.Text()] = true
		f.Service.Save(sc.Text())
	}
	// f.mem = NewInMemoryMap()
	// f.mem.cache = initCached
	// log.Info("load data from file done, file = ", f.File, " len = ", len(initCached))

	b := gobatch.NewBatcher(f.Flush)
	b.MaxRetry = f.Retry
	b.BatchSize = f.Batch
	b.MaxWait, _ = time.ParseDuration(f.Interval)

	ch, err := b.Start(context.TODO())
	if err != nil {
		return err
	}
	f.ch = ch

	return nil
}

func (f *PersistedIdempotent) Duplicated(key interface{}) (bool, error) {
	return f.Service.Duplicated(key)
}
func (f *PersistedIdempotent) Save(key interface{}) error {
	err := f.Service.Save(key)
	if err != nil {
		return err
	}
	f.ch <- key

	return nil
}

func (f *PersistedIdempotent) AllKeys() ([]interface{}, error) {
	return f.Service.AllKeys()
}

func (f *PersistedIdempotent) Flush(ctx context.Context, queue []interface{}) error {
	file, err := os.OpenFile(f.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, item := range queue {
		txt := fmt.Sprintf("%s\n", item)
		file.WriteString(txt)
	}
	log.Info("append to file done, len ", len(queue))

	return nil
}
