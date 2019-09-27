package errors

import "reflect"

var _errTags = make(map[string]bool)

// ErrTagRegistry errors
func ErrTagRegistry(tags ...interface{}) {
	for _, tag := range tags {
		if tag == nil || _isNone(tag) {
			continue
		}

		var _tags []string
		t := reflect.ValueOf(tag)
		switch t.Kind() {
		case reflect.Ptr, reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				_tags = append(_tags, t.Field(i).String())
			}
		case reflect.String:
			_tags = append(_tags, tag.(string))
		}

		for _, t := range _tags {
			if _, ok := _errTags[t]; ok {
				PanicT(ok, "tag %s has existed", t)
			}
			_errTags[t] = true
		}
	}
}

// ErrTags get error tags
func ErrTags() map[string]bool {
	return _errTags
}

// ErrTagsMatch check error tags
func ErrTagsMatch(tag string) bool {
	_, ok := _errTags[tag]
	return ok
}
