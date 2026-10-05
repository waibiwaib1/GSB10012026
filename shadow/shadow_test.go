package shadow

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	mqtt "github.com/at-wat/mqtt-go"
	mockmqtt "github.com/at-wat/mqtt-go/mock"
)

type mockDevice struct {
	*mockmqtt.Client
}

func (d *mockDevice) ThingName() string {
	return "test"
}

func TestShadow(t *testing.T) {
	t.Run("Rejected", func(t *testing.T) {
		chErr := make(chan error, 1)
		s := &shadow{
			doc: &ThingDocument{
				State: ThingState{
					Desired: map[string]interface{}{"key": "value_init"},
				},
				Version:   2,
				Timestamp: 12345,
			},
			onDelta: func(delta map[string]interface{}) {
				t.Error("onDelta must not be called on rejected")
			},
			onError: func(err error) {
				chErr <- err
			},
		}
		expectedDoc := &ThingDocument{
			State: ThingState{
				Desired: map[string]interface{}{"key": "value_init"},
			},
			Version:   2,
			Timestamp: 12345,
		}

		t.Run("Get", func(t *testing.T) {
			s.getRejected(&mqtt.Message{})
			select {
			case err := <-chErr:
				if err == nil {
					t.Error("onError must be called with non-nil error")
				}
			default:
				t.Fatal("Timeout")
			}
			if !reflect.DeepEqual(*expectedDoc, *s.doc) {
				t.Error("Document must not be changed on reject")
			}
		})
		t.Run("Update", func(t *testing.T) {
			s.updateRejected(&mqtt.Message{})
			select {
			case err := <-chErr:
				if err == nil {
					t.Error("onError must be called with non-nil error")
				}
			default:
				t.Fatal("Timeout")
			}
			if !reflect.DeepEqual(*expectedDoc, *s.doc) {
				t.Error("Document must not be changed on reject")
			}
		})
		t.Run("Delete", func(t *testing.T) {
			s.deleteRejected(&mqtt.Message{})
			select {
			case err := <-chErr:
				if err == nil {
					t.Error("onError must be called with non-nil error")
				}
			default:
				t.Fatal("Timeout")
			}
			if !reflect.DeepEqual(*expectedDoc, *s.doc) {
				t.Error("Document must not be changed on reject")
			}
		})
	})

	t.Run("Accepted", func(t *testing.T) {
		t.Run("Get", func(t *testing.T) {
			chErr := make(chan error, 1)
			chDelta := make(chan map[string]interface{}, 1)
			s := &shadow{
				doc: &ThingDocument{
					State: ThingState{
						Desired: map[string]interface{}{"key": "value_init"},
					},
					Version:   2,
					Timestamp: 12345,
				},
				onDelta: func(delta map[string]interface{}) {
					chDelta <- delta
				},
				onError: func(err error) {
					chErr <- err
				},
			}
			expectedDoc := &ThingDocument{
				State: ThingState{
					Desired: map[string]interface{}{"key2": "value2"},
					Delta:   map[string]interface{}{"key2": "value2"},
				},
				Version:   3,
				Timestamp: 12346,
			}
			s.getAccepted(&mqtt.Message{Payload: []byte(
				"{" +
					"  \"state\": {" +
					"    \"desired\": {\"key2\": \"value2\"}," +
					"    \"delta\": {\"key2\": \"value2\"}" +
					"  }," +
					"  \"version\": 3," +
					"  \"timestamp\": 12346" +
					"}",
			)})

			select {
			case err := <-chErr:
				t.Error(err)
			case delta := <-chDelta:
				if !reflect.DeepEqual(expectedDoc.State.Delta, delta) {
					t.Errorf("Expected delta: %v, got: %v",
						expectedDoc.State.Delta, delta,
					)
				}
			default:
				t.Fatal("Timeout")
			}
			if !reflect.DeepEqual(expectedDoc, s.doc) {
				t.Errorf("Expected state: %v, got: %v",
					expectedDoc, s.doc,
				)
			}
		})
		t.Run("Delete", func(t *testing.T) {
			chErr := make(chan error, 1)
			s := &shadow{
				doc: &ThingDocument{
					State: ThingState{
						Desired: map[string]interface{}{"key": "value_init"},
					},
					Version:   2,
					Timestamp: 12345,
				},
				onDelta: func(delta map[string]interface{}) {
					t.Error("Delete must not trigger onDelta")
				},
				onError: func(err error) {
					chErr <- err
				},
			}
			s.deleteAccepted(&mqtt.Message{})

			select {
			case err := <-chErr:
				t.Error(err)
			default:
			}
			if s.doc != nil {
				t.Errorf("Document must be nil after delete")
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Delta", func(t *testing.T) {
			chErr := make(chan error, 1)
			chDelta := make(chan map[string]interface{}, 1)
			s := &shadow{
				doc: &ThingDocument{
					State: ThingState{
						Desired: map[string]interface{}{"key2": "value2"},
					},
					Version:   2,
					Timestamp: 12345,
				},
				onDelta: func(delta map[string]interface{}) {
					chDelta <- delta
				},
				onError: func(err error) {
					chErr <- err
				},
			}
			expectedDoc := &ThingDocument{
				State: ThingState{
					Desired: map[string]interface{}{"key2": "value2"},
					Delta:   map[string]interface{}{"key2": "value2"},
				},
				Version:   3,
				Timestamp: 12346,
			}
			s.updateDelta(&mqtt.Message{Payload: []byte(
				"{" +
					"  \"state\": {" +
					"    \"key2\": \"value2\"" +
					"  }," +
					"  \"version\": 3," +
					"  \"timestamp\": 12346" +
					"}",
			)})

			select {
			case err := <-chErr:
				t.Error(err)
			case delta := <-chDelta:
				if !reflect.DeepEqual(expectedDoc.State.Delta, delta) {
					t.Errorf("Expected delta: %v, got: %v",
						expectedDoc.State.Delta, delta,
					)
				}
			default:
				t.Fatal("Timeout")
			}
			if !reflect.DeepEqual(expectedDoc, s.doc) {
				t.Errorf("Expected state: %v, got: %v",
					expectedDoc, s.doc,
				)
			}
		})
		t.Run("Accepted", func(t *testing.T) {
			chErr := make(chan error, 1)
			chDelta := make(chan map[string]interface{}, 1)
			s := &shadow{
				doc: &ThingDocument{
					State: ThingState{
						Reported: map[string]interface{}{"key1": "value1"},
					},
					Version:   2,
					Timestamp: 12345,
				},
				onDelta: func(delta map[string]interface{}) {
					chDelta <- delta
				},
				onError: func(err error) {
					chErr <- err
				},
			}
			expectedDoc := &ThingDocument{
				State: ThingState{
					Reported: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Version:   3,
				Timestamp: 12346,
			}
			s.updateAccepted(&mqtt.Message{Payload: []byte(
				"{" +
					"  \"state\": {" +
					"    \"Reported\": {" +
					"      \"key2\": \"value2\"" +
					"    }" +
					"  }," +
					"  \"version\": 3," +
					"  \"timestamp\": 12346" +
					"}",
			)})

			select {
			case err := <-chErr:
				t.Error(err)
			case <-chDelta:
				t.Error("cbDelta must not be called on update accept")
			default:
			}
			if !reflect.DeepEqual(expectedDoc, s.doc) {
				t.Errorf("Expected state: %v, got: %v",
					expectedDoc, s.doc,
				)
			}
		})
		t.Run("OldDelta", func(t *testing.T) {
			s := &shadow{
				doc: &ThingDocument{
					State: ThingState{
						Desired: map[string]interface{}{"key2": "value2"},
					},
					Version:   2,
					Timestamp: 12345,
				},
				onDelta: func(delta map[string]interface{}) {
					t.Error("Old delta must be discarded")
				},
				onError: func(err error) {
					t.Error("Old delta must be silently discarded")
				},
			}
			expectedDoc := &ThingDocument{
				State: ThingState{
					Desired: map[string]interface{}{"key2": "value2"},
				},
				Version:   2,
				Timestamp: 12345,
			}
			s.updateDelta(&mqtt.Message{Payload: []byte(
				"{" +
					"  \"state\": {" +
					"    \"key2\": \"value\"" +
					"  }," +
					"  \"version\": 1," +
					"  \"timestamp\": 12343" +
					"}",
			)})

			if !reflect.DeepEqual(expectedDoc, s.doc) {
				t.Errorf("Expected state: %v, got: %v",
					expectedDoc, s.doc,
				)
			}
		})
	})
}

func TestClientToken(t *testing.T) {
	tests := map[string]struct {
		call          func(context.Context, Shadow) error
		requestTopic  string
		responseTopic string
		rejected      bool
	}{
		"ReportAccepted": {
			call: func(ctx context.Context, s Shadow) error {
				return s.Report(ctx, map[string]interface{}{"key": "value"})
			},
			requestTopic:  "update",
			responseTopic: "update/accepted",
		},
		"ReportRejected": {
			call: func(ctx context.Context, s Shadow) error {
				return s.Report(ctx, map[string]interface{}{"key": "value"})
			},
			requestTopic:  "update",
			responseTopic: "update/rejected",
			rejected:      true,
		},
		"DesireAccepted": {
			call: func(ctx context.Context, s Shadow) error {
				return s.Desire(ctx, map[string]interface{}{"key": "value"})
			},
			requestTopic:  "update",
			responseTopic: "update/accepted",
		},
		"DesireRejected": {
			call: func(ctx context.Context, s Shadow) error {
				return s.Desire(ctx, map[string]interface{}{"key": "value"})
			},
			requestTopic:  "update",
			responseTopic: "update/rejected",
			rejected:      true,
		},
		"DeleteAccepted": {
			call:          func(ctx context.Context, s Shadow) error { return s.Delete(ctx) },
			requestTopic:  "delete",
			responseTopic: "delete/accepted",
		},
		"DeleteRejected": {
			call:          func(ctx context.Context, s Shadow) error { return s.Delete(ctx) },
			requestTopic:  "delete",
			responseTopic: "delete/rejected",
			rejected:      true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			var sh Shadow
			cli := &mockDevice{
				Client: &mockmqtt.Client{
					PublishFn: func(ctx context.Context, msg *mqtt.Message) error {
						if msg.Topic != sh.(*shadow).topic(test.requestTopic) {
							t.Errorf("Expected topic: %s, got: %s", test.requestTopic, msg.Topic)
						}
						req := &thingDocumentRaw{}
						if test.requestTopic == "delete" {
							deleteReq := &simpleRequest{}
							if err := json.Unmarshal(msg.Payload, deleteReq); err != nil {
								return err
							}
							if deleteReq.ClientToken == "" {
								t.Fatal("clientToken must be set")
							}
							serveShadowResponse(t, sh, test, deleteReq.ClientToken)
							return nil
						}
						if err := json.Unmarshal(msg.Payload, req); err != nil {
							return err
						}
						if req.ClientToken == "" {
							t.Fatal("clientToken must be set")
						}
						serveShadowResponse(t, sh, test, req.ClientToken)
						return nil
					},
				},
			}

			var err error
			sh, err = New(ctx, cli)
			if err != nil {
				t.Fatal(err)
			}
			cli.Handle(sh)

			err = test.call(ctx, sh)
			if test.rejected {
				if err == nil {
					t.Fatal("expected rejection error")
				}
				errResp, ok := err.(*ErrorResponse)
				if !ok {
					t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
				}
				if errResp.Code != 400 || errResp.Message != "Bad Request" {
					t.Fatalf("unexpected error response: %v", errResp)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func serveShadowResponse(t *testing.T, sh Shadow, test struct {
	call          func(context.Context, Shadow) error
	requestTopic  string
	responseTopic string
	rejected      bool
}, token string) {
	t.Helper()

	var payload []byte
	if test.rejected {
		payload, _ = json.Marshal(&ErrorResponse{
			Code:        400,
			Message:     "Bad Request",
			ClientToken: token,
		})
	} else if test.requestTopic == "delete" {
		payload, _ = json.Marshal(&deleteResponse{
			Version:     1,
			ClientToken: token,
		})
	} else {
		payload, _ = json.Marshal(&thingDocumentRaw{
			State: thingStateRaw{
				Reported: json.RawMessage(`{"key":"value"}`),
			},
			Version:     1,
			ClientToken: token,
		})
	}

	sh.Serve(&mqtt.Message{
		Topic:   sh.(*shadow).topic(test.responseTopic),
		Payload: payload,
	})
}
