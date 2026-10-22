package test

import (
	"reflect"
	"testing"

	"github.com/gin-contrib/sessions"

	"opennamu/route"
	"opennamu/route/tool"
)

// memorySession is an in-memory sessions.Session used to verify that
// View_w_add_recent_document stores and retrieves lastest_document correctly.
// It follows the same sessions.Session interface as the existing emptySession
// but provides real get/set/save so round-trip behavior can be checked.
type memorySession struct {
	values map[interface{}]interface{}
}

func New_memory_session() *memorySession {
	return &memorySession{values: map[interface{}]interface{}{}}
}

func (m *memorySession) ID() string { return "" }

func (m *memorySession) Get(key interface{}) interface{} {
	value, ok := m.values[key]
	if !ok {
		return nil
	}
	return value
}

func (m *memorySession) Set(key, value interface{}) {
	m.values[key] = value
}

func (m *memorySession) Delete(key interface{}) {
	delete(m.values, key)
}

func (m *memorySession) Clear() {
	m.values = map[interface{}]interface{}{}
}

func (m *memorySession) AddFlash(value interface{}, vars ...string) {}

func (m *memorySession) Flashes(vars ...string) []interface{} { return nil }

func (m *memorySession) Options(o sessions.Options) {}

func (m *memorySession) Save() error { return nil }

// recent_docs returns the documents stored in the session.
func (m *memorySession) recent_docs() []string {
	value, ok := m.values["lastest_document"].([]string)
	if !ok {
		return []string{}
	}
	return append([]string{}, value...)
}

func Test_view_w_add_recent_document(t *testing.T) {
	// initial documents in oldest-first order; nil means empty session.
	mk := func(initial []string) (tool.Config, *memorySession) {
		session := New_memory_session()
		if initial != nil {
			session.values["lastest_document"] = initial
		}
		return tool.Config{IP: "test_ip", Session: session}, session
	}

	// check verifies the returned list and that the session stores the same
	// value (return value must equal stored value).
	check := func(name string, initial []string, doc string, want []string) {
		config, session := mk(initial)
		got := route.View_w_add_recent_document(config, doc)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", name, got, want)
		}
		stored := session.recent_docs()
		if !reflect.DeepEqual(stored, want) {
			t.Errorf("%s: session stored %v, want %v", name, stored, want)
		}
	}

	// 1. empty session -> [A]
	check("empty add A", nil, "A", []string{"A"})

	// 2. 9 records (A-I), add J -> 10
	check("nine add J", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
		"J", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"})

	// 3. 10 records (A-J), J re-visit -> 10, no dup, oldest kept
	check("ten J revisit", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
		"J", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"})

	// 4. 10 records (A-J), E re-visit -> moves E to end
	check("ten E revisit", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
		"E", []string{"A", "B", "C", "D", "F", "G", "H", "I", "J", "E"})

	// 5. 10 records (A-J), A re-visit -> moves A to end
	check("ten A revisit", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
		"A", []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "A"})

	// 6. 10 records (A-J), new K -> oldest A dropped
	check("ten new K", []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
		"K", []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K"})

	// 7. consecutive same doc -> no dup, no unnecessary delete
	check("two A revisit", []string{"A", "B"}, "A", []string{"B", "A"})
	check("single A revisit", []string{"A"}, "A", []string{"A"})

	// 8a. empty doc name -> original list returned and unchanged
	check("empty doc name", []string{"A", "B", "C"}, "",
		[]string{"A", "B", "C"})

	// 8b. nil session -> empty list
	nil_session_config := tool.Config{IP: "test_ip"}
	nil_session_got := route.View_w_add_recent_document(nil_session_config, "A")
	if !reflect.DeepEqual(nil_session_got, []string{}) {
		t.Errorf("nil session: got %v, want []", nil_session_got)
	}

	// 9. session stored value equals return value for a representative case
	config, session := mk([]string{"A", "B", "C"})
	got := route.View_w_add_recent_document(config, "D")
	if !reflect.DeepEqual(got, session.recent_docs()) {
		t.Errorf("session mismatch: got %v, session %v", got, session.recent_docs())
	}
	if !reflect.DeepEqual(got, []string{"A", "B", "C", "D"}) {
		t.Errorf("representative case: got %v, want [A B C D]", got)
	}
}

func Test_view_w_recent_documents(t *testing.T) {
	// Returns stored list (a copy) and handles nil session / missing key.
	check := func(initial []string, want []string) {
		config := tool.Config{IP: "test_ip"}
		if initial != nil {
			session := New_memory_session()
			session.values["lastest_document"] = initial
			config.Session = session
		}
		got := route.View_w_recent_documents(config)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("View_w_recent_documents: got %v, want %v", got, want)
		}
	}

	check([]string{}, []string{})
	check(nil, []string{})
	check([]string{"A", "B"}, []string{"A", "B"})
}
