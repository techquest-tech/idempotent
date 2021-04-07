package idempotent

import (
	"bytes"
	"text/template"
	"time"
)

type DefaultIdempotentKey struct {
	Target      interface{}
	KeysTmpl    string
	IgnoreError bool
	tmpl        *template.Template
}

func FormateDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

func NewIdempotent(keys string) (*DefaultIdempotentKey, error) {
	out := &DefaultIdempotentKey{
		KeysTmpl: keys,
	}

	tp, err := template.New("idempotent").Funcs(template.FuncMap{
		"FormateDate": FormateDate,
	}).Parse(keys)

	if err != nil {
		// log.Fatal("key template error, ", err)
		log.Error("key template error, ", err)
		return nil, err
	}

	out.tmpl = tp
	return out, nil
}

func (d DefaultIdempotentKey) IdempotentKey() (interface{}, error) {

	out := bytes.Buffer{}
	err := d.tmpl.Execute(&out, d.Target)
	if err != nil {
		log.Error("exec template error, ", err)
		return nil, err
	}

	return out.String(), nil
}
