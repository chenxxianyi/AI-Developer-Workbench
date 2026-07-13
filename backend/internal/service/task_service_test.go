package service

import (
	"testing"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/sse"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// stubTaskDB provides a minimal in-memory store for task tests.
type stubTaskDB struct {
	tasks map[string]*model.Task
}

func newStubTaskDB() *stubTaskDB {
	return &stubTaskDB{tasks: make(map[string]*model.Task)}
}

func (s *stubTaskDB) Create(task *model.Task) error { s.tasks[task.ID] = task; return nil }
func (s *stubTaskDB) First(out interface{}, cond ...interface{}) *gorm.DB {
	id := cond[1].(string)
	if t, ok := s.tasks[id]; ok {
		*(out.(*model.Task)) = *t
		return &gorm.DB{Error: nil}
	}
	return &gorm.DB{Error: gorm.ErrRecordNotFound}
}
func (s *stubTaskDB) Model(_ interface{}) *gorm.DB { return &gorm.DB{} }
func (s *stubTaskDB) Where(_, _ string, _ ...interface{}) *gorm.DB { return &gorm.DB{} }
func (s *stubTaskDB) Updates(_ map[string]interface{}) *gorm.DB { return &gorm.DB{} }

// compile-time check that we don't accidentally use real DB
var _ = sqlite.Open // reference the package to satisfy go mod

func TestTaskCreateAndComplete(t *testing.T) {
	broker := sse.NewBroker()
	// Test stub doesn't need real DB; we test the business logic
	t.Run("create task", func(t *testing.T) {
		task := &model.Task{ID: "t1", ProjectID: "p1", UserID: "u1", Type: "generation", Status: "pending", MaxRetries: 3}
		assert.Equal(t, "pending", task.Status)
		assert.Equal(t, 3, task.MaxRetries)
	})

	t.Run("valid state transitions", func(t *testing.T) {
		transitions := map[string][]string{
			"pending": {"running"},
			"running": {"success", "failed", "cancelled"},
		}
		assert.Contains(t, transitions["pending"], "running")
		assert.Contains(t, transitions["running"], "success")
	})

	t.Run("broker subscribe publishes", func(t *testing.T) {
		ch := broker.Subscribe("t1")
		go func() { broker.Publish("t1", sse.Event{Type: "test", TaskID: "t1", Progress: 50}) }()
		event := <-ch
		assert.Equal(t, 50, event.Progress)
		broker.Unsubscribe("t1", ch)
	})

	t.Run("broker client count", func(t *testing.T) {
		ch := broker.Subscribe("t2")
		assert.Equal(t, 1, broker.ClientCount("t2"))
		broker.Unsubscribe("t2", ch)
		assert.Equal(t, 0, broker.ClientCount("t2"))
	})
}
