//go:build integration

package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/joaopandolfi/blackwhale/v2/remotes/request"
	"github.com/tj/assert"
)

func TestA(t *testing.T) {
	p, _ := json.Marshal(map[string]any{
		"data": map[string]any{
			"created_at": "2024-01-15T15:03:33.165626Z",
			"delivery_info": map[string]any{
				"current_retry": 324,
				"max_retries":   3,
			},
			"event": map[string]any{
				"data": map[string]any{
					"new": map[string]any{
						"created_at":       "2024-01-15T15:03:33.165686+00:00",
						"description":      "Justificativa muito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muitomuito muito Grande",
						"entity":           "agreement",
						"entity_id":        "e8fe61ca-5db8-4e78-a02a-3cff5757b74b",
						"finished_at":      nil,
						"ids":              "e43fb40d-3f3c-4de8-bc6d-4a8eede6f5d4",
						"participation_id": "fe9f04c0-d089-4998-9d2a-e734c8c35376",
						"requester_id":     "8440deec-f752-46b2-95c1-fe9faba4cdbf",
						"resolver_id":      "8440deec-f752s-46b2-95c1-fe9faba4cdbf",
						"status":           "pending",
						"updated_at":       "2024-01-15T15:03:33.165686+00:00",
					},
				},
				"session_variables": map[string]any{
					"x-hasura-role":    "system",
					"x-hasura-user-id": "sauron",
				},
				"trace_context": map[string]any{
					"span_id":  "53ccd0370875a7cf",
					"trace_id": "c2d175a9f1e00dd3cfe1f966f19b8721",
				},
			},
			"updated_by": "cf2f1e74-eb77-4e17-8831-093495148590",
		},
		"emitter":        "hasura-event-trigger-new-adjustment-request",
		"name":           "Hasura Event Trigger New Adjustment Request",
		"retry_ttl":      2,
		"subtype":        "participants",
		"system_destiny": "logger",
		"system_origin":  "hasura",
		"type":           "adjustment_request",
	},
	)
	r, code, _ := request.PostWithHeader2("http://localhost:5656/event/new", map[string]string{
		"Content-Encoding": "gzip",
	}, p)

	assert.Nil(t, string(r))
	assert.Equal(t, code, 200)
}
