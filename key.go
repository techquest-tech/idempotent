package idempotent

import (
	"bytes"
	"text/template"
	"time"

	"go.uber.org/zap"
)

type DefaultIdempotentKey struct {
	Target      interface{}
	Keys        FuncKeys
	KeysTmpl    string
	IgnoreError bool
	tmpl        *template.Template
	Logger      *zap.Logger
}

type FuncKeys func(obj interface{}) (interface{}, error)

func FormateDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

// use golang tempate as key value.
func TemplateAsKey(keys string) (*DefaultIdempotentKey, error) {
	out := &DefaultIdempotentKey{
		KeysTmpl: keys,
		// Logger:   logger,
	}

	tp, err := template.New("idempotent").Funcs(template.FuncMap{
		"FormateDate": FormateDate,
	}).Parse(keys)

	if err != nil {
		// log.Fatal("key template error, ", err)
		// logger.Error("key template error, ", zap.Any("error", err))
		return nil, err
	}

	out.tmpl = tp
	return out, nil
}

func (d DefaultIdempotentKey) IdempotentKey() (interface{}, error) {
	if d.Keys != nil {
		return d.Keys(d.Target)
	}

	out := bytes.Buffer{}
	err := d.tmpl.Execute(&out, d.Target)
	if err != nil {
		d.Logger.Error("exec template error, ", zap.Any("error", err))
		return nil, err
	}

	return out.String(), nil
}
