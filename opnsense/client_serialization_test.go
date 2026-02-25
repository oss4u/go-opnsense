package opnsense

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpnSenseApi_SerializesConcurrentRequests(t *testing.T) {
	var inFlight int32
	var maxInFlight int32

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || user != "test-key" || pass != "test-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		current := atomic.AddInt32(&inFlight, 1)
		for {
			seen := atomic.LoadInt32(&maxInFlight)
			if current <= seen {
				break
			}
			if atomic.CompareAndSwapInt32(&maxInFlight, seen, current) {
				break
			}
		}

		time.Sleep(40 * time.Millisecond)

		atomic.AddInt32(&inFlight, -1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	api := NewOpnSenseClient(server.URL, "test-key", "test-secret")

	var wg sync.WaitGroup
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func(idx int) {
			defer wg.Done()
			_, err := api.ModifyingRequest("core", "service", "set", `{"idx":1}`, []string{"svc"})
			require.NoError(t, err, "request %d failed", idx)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int32(1), atomic.LoadInt32(&maxInFlight), "only one request should be in-flight at a time")
}
