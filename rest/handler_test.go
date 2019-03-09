package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/xid"
)

func TestPOSTHandler(t *testing.T) {

	tests := []struct {
		name             string
		withPayloadError bool
		withCreateError  bool
		wantedStatus     int
	}{
		{
			name:         "working POST",
			wantedStatus: http.StatusCreated,
		},
		{
			name:             "bad payload, non working POST",
			withPayloadError: true,
			wantedStatus:     http.StatusBadRequest,
		},
		{
			name:            "internal error, non working POST",
			withCreateError: true,
			wantedStatus:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			e.StatusID = 1
			payload, _ := json.Marshal(e)
			if tt.withPayloadError {
				payload = payload[1:]
			}
			m.wantCreateError = tt.withCreateError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST", `http://dummy/entity`,
				bytes.NewBuffer(payload),
			)

			POSTHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
				return
			}

			if resp.StatusCode == http.StatusOK {
				entx := &entityMock{}
				_ = json.NewDecoder(resp.Body).Decode(entx)
				if entx.StatusID != 1 {
					t.Errorf("POSTHandler wasn't created properly")
				}
			}
		})
	}
}

func BenchmarkPOSTHandler(b *testing.B) {
	m := newMgrMock()
	payload, _ := json.Marshal(&entityMock{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rr := httptest.NewRecorder()
		jsonRequestOK := httptest.NewRequest("POST", `http://dummy/ent`,
			bytes.NewBuffer(payload))
		b.StartTimer()
		POSTHandler(m)(rr, jsonRequestOK)
	}
}
func TestGETListHandler(t *testing.T) {

	tests := []struct {
		name             string
		withEmptyList    bool
		withGetListError bool
		wantedStatus     int
	}{
		{
			name:         "working GETList",
			wantedStatus: http.StatusOK,
		},
		{
			name:          "working empty GETList",
			withEmptyList: true,
			wantedStatus:  http.StatusNotFound,
		},
		{
			name:             "internal error, non working GETList",
			withGetListError: true,
			wantedStatus:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			if !tt.withEmptyList {
				for i := 0; i < 5; i++ {
					_, _ = m.Create(ctx, e, bytes.NewReader([]byte{}))
				}
			}
			m.wantGetListError = tt.withGetListError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET", `http://dummy/entity?pol=lux`, bytes.NewReader([]byte{}),
			)

			GETListHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
				return
			}

			if resp.StatusCode == http.StatusOK {
				type em []*entityMock
				entx := em{}
				if err := json.NewDecoder(resp.Body).Decode(&entx); err != nil {
					t.Errorf("GETListHandler: %v", err)
					return
				}
				if len(entx) != 5 {
					t.Errorf("GETListHandler: got %d ent instead of %d",
						len(entx), 5)
				}
			}
		})
	}
}

func BenchmarkGETListHandler(b *testing.B) {
	ctx := context.Background()
	m := newMgrMock()
	e := m.NewEmptyEntity().(*entityMock)
	for i := 0; i < 5; i++ {
		_, _ = m.Create(ctx, e, bytes.NewReader([]byte{}))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET", `http://dummy/entity`, bytes.NewReader([]byte{}),
		)
		b.StartTimer()
		GETListHandler(m)(rr, req)
	}
}

func TestGETHandler(t *testing.T) {

	tests := []struct {
		name         string
		withEmpty    bool
		withGetError bool
		withBadID    bool
		wantedStatus int
	}{
		{
			name:         "working GET",
			wantedStatus: http.StatusOK,
		},
		{
			name:         "working empty GET",
			withEmpty:    true,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "internal error, non working GET",
			withGetError: true,
			wantedStatus: http.StatusInternalServerError,
		},
		{
			name:         "id error, non working GET",
			withBadID:    true,
			wantedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			ecc := &entityMock{}
			reqID := xid.New()
			if !tt.withEmpty {
				ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
				ecc = ec.(*entityMock)
				reqID = ecc.ID
			}
			if tt.withBadID {
				reqID = xid.ID{}
			}
			m.wantGetError = tt.withGetError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET", `http://dummy/entity`, bytes.NewReader([]byte{}),
			)
			req = req.WithContext(GetTestContextWithID(req.Context(), reqID))

			GETHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf(
					"GETHandler: returned wrong status code: got %v want %v",
					status, tt.wantedStatus,
				)
				return
			}

			if resp.StatusCode == http.StatusOK {
				entx := entityMock{}
				if err := json.NewDecoder(resp.Body).Decode(&entx); err != nil {
					t.Errorf("GETHandler: %v", err)
					return
				}
				if entx.ID != ecc.ID {
					t.Errorf("GETHandler: didn't get the right entity")
				}
			}
		})
	}
}

func BenchmarkGETHandler(b *testing.B) {
	ctx := context.Background()
	m := newMgrMock()
	e := m.NewEmptyEntity().(*entityMock)
	ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
	ecc := ec.(*entityMock)
	reqID := ecc.ID
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET", `http://dummy/entity`, bytes.NewReader([]byte{}),
		)
		req = req.WithContext(GetTestContextWithID(req.Context(), reqID))
		b.StartTimer()
		GETHandler(m)(rr, req)
	}
}

func TestDELETEHandler(t *testing.T) {

	tests := []struct {
		name            string
		withEmpty       bool
		withDeleteError bool
		withBadID       bool
		wantedStatus    int
	}{
		{
			name:         "working DELETE",
			wantedStatus: http.StatusAccepted,
		},
		{
			name:         "working empty DELETE",
			withEmpty:    true,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:            "internal error, non working DELETE",
			withDeleteError: true,
			wantedStatus:    http.StatusInternalServerError,
		},
		{
			name:         "id error, non working DELETE",
			withBadID:    true,
			wantedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			reqID := xid.New()
			if !tt.withEmpty {
				ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
				ecc := ec.(*entityMock)
				reqID = ecc.ID
			}
			if tt.withBadID {
				reqID = xid.ID{}
			}
			m.wantDeleteError = tt.withDeleteError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"DELETE", `http://dummy/entity`, bytes.NewReader([]byte{}),
			)
			req = req.WithContext(GetTestContextWithID(req.Context(), reqID))

			DELETEHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf(
					"DELETEHandler: returned wrong status code: got %v want %v",
					status, tt.wantedStatus,
				)
				return
			}

			if resp.StatusCode == http.StatusAccepted {
				if _, ok := m.EntityList[reqID]; ok {
					t.Errorf("DELETEHandler: didn't delete the entity")
				}
			}
		})
	}
}

func BenchmarkDELETEHandler(b *testing.B) {
	ctx := context.Background()
	m := newMgrMock()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		e := m.NewEmptyEntity().(*entityMock)
		ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
		ecc := ec.(*entityMock)
		reqID := ecc.ID
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"DELETE", `http://dummy/entity`, bytes.NewReader([]byte{}),
		)
		req = req.WithContext(GetTestContextWithID(req.Context(), reqID))
		b.StartTimer()
		DELETEHandler(m)(rr, req)
	}
}

func TestPUTHandler(t *testing.T) {

	tests := []struct {
		name             string
		withPayloadError bool
		withEmpty        bool
		withUpdateError  bool
		withBadID        bool
		wantedStatus     int
	}{
		{
			name:         "working PUT",
			wantedStatus: http.StatusOK,
		},
		{
			name:             "bad payload, non working PUT",
			withPayloadError: true,
			wantedStatus:     http.StatusBadRequest,
		},
		{
			name:         "non existing id, non working PUT",
			withEmpty:    true,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "bad id, non working PUT",
			withBadID:    true,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:            "internal error, non working PUT",
			withUpdateError: true,
			wantedStatus:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			reqID := xid.New()
			if !tt.withEmpty {
				ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
				ecc := ec.(*entityMock)
				reqID = ecc.ID
			}
			if tt.withBadID {
				reqID = xid.ID{}
			}
			e.StatusID = 2
			payload, _ := json.Marshal(e)
			if tt.withPayloadError {
				payload = payload[1:]
			}
			m.wantUpdateError = tt.withUpdateError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT", `http://dummy/entity`,
				bytes.NewBuffer(payload),
			)
			req = req.WithContext(GetTestContextWithID(req.Context(), reqID))

			PUTHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
				return
			}

			if resp.StatusCode == http.StatusOK {
				entx := &entityMock{}
				_ = json.NewDecoder(resp.Body).Decode(entx)
				if entx.StatusID != 2 {
					t.Errorf("PUTHandler didn't update properly")
				}
			}
		})
	}
}

func BenchmarkPUTHandler(b *testing.B) {
	ctx := context.Background()
	m := newMgrMock()
	e := m.NewEmptyEntity().(*entityMock)
	ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
	ecc := ec.(*entityMock)
	reqID := ecc.ID
	e.StatusID = 2
	payload, _ := json.Marshal(e)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", `http://dummy/ent`,
			bytes.NewBuffer(payload))
		req = req.WithContext(GetTestContextWithID(req.Context(), reqID))
		b.StartTimer()
		PUTHandler(m)(rr, req)
	}
}

func TestPATCHHandler(t *testing.T) {

	tests := []struct {
		name                   string
		withPayloadError       bool
		withEmpty              bool
		withPartialUpdateError bool
		withBadID              bool
		wantedStatus           int
	}{
		{
			name:         "working PATCH",
			wantedStatus: http.StatusNoContent,
		},
		{
			name:             "bad payload, non working PATCH",
			withPayloadError: true,
			wantedStatus:     http.StatusBadRequest,
		},
		{
			name:         "non existing id, non working PATCH",
			withEmpty:    true,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "bad id, non working PATCH",
			withBadID:    true,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:                   "internal error, non working PATCH",
			withPartialUpdateError: true,
			wantedStatus:           http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := newMgrMock()
			e := m.NewEmptyEntity().(*entityMock)
			reqID := xid.New()
			if !tt.withEmpty {
				ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
				ecc := ec.(*entityMock)
				reqID = ecc.ID
			}
			if tt.withBadID {
				reqID = xid.ID{}
			}
			payload := []byte(`{"status_id": 2}`)
			if tt.withPayloadError {
				payload = payload[1:]
			}
			m.wantPartialUpdateError = tt.withPartialUpdateError

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PATCH", `http://dummy/entity`,
				bytes.NewBuffer(payload),
			)
			req = req.WithContext(GetTestContextWithID(req.Context(), reqID))

			PATCHHandler(m)(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
				return
			}

			if resp.StatusCode == http.StatusOK {
				entx := &entityMock{}
				_ = json.NewDecoder(resp.Body).Decode(entx)
				if entx.StatusID != 2 {
					t.Errorf("PATCHHandler didn't update properly")
				}
			}
		})
	}
}

func BenchmarkPATCHHandler(b *testing.B) {
	ctx := context.Background()
	m := newMgrMock()
	e := m.NewEmptyEntity().(*entityMock)
	ec, _ := m.Create(ctx, e, bytes.NewReader([]byte{}))
	ecc := ec.(*entityMock)
	reqID := ecc.ID
	payload := []byte(`{"status_id": 2}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("PATCH", `http://dummy/ent`,
			bytes.NewBuffer(payload))
		req = req.WithContext(GetTestContextWithID(req.Context(), reqID))
		b.StartTimer()
		PATCHHandler(m)(rr, req)
	}
}
