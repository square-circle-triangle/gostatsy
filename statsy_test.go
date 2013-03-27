package statsy

import (
	"errors"
	"testing"
  "time"
)

func NewTestClient() *Statsy {
	statsyClient := New("abc123", "def456")
  statsyClient.BaseUrl = testHttpServer.BaseURL()
  statsyClient.expires = time.Date(2023, 3, 27, 13, 35, 53, 0, time.UTC).Unix()
  return statsyClient
}

var testHttpServer = NewHTTPServer()

func TestCreateNewClient(t *testing.T) {
	client := NewTestClient()

	if client.ApiKey != "abc123" {
		t.Error("API should be abc123 but got", client.ApiKey)
	}

	if client.SecretKey != "def456" {
		t.Error("SecretKey should be def456 but got", client.SecretKey)
	}
}

func TestSignatureCalculationWithNoStreamPrefix(t *testing.T) {
	client := NewTestClient()

	signature := client.Sign(300, "")
	correctSignature := "6cpng1sPY9htddzYA/2uADMRovA="

	if signature != correctSignature {
		t.Error("Expected signature to be", correctSignature, "but got", signature)
	}
}

func TestSignatureCalculationWithStreamPrefix(t *testing.T) {
	client := NewTestClient()

	signature := client.Sign(600, "widgets.melbourne")
	correctSignature := "R1Ro0z/gQw5Mu/R+aKQp7tL7sKE="

	if signature != correctSignature {
		t.Error("Expected signature to be", correctSignature, "but got", signature)
	}
}

func TestSendSingleEvent(t *testing.T) {
	client := NewTestClient()
	events := make([]Event, 0)

	events = append(events, Event{Stream: "haproxy.requests"})
	client.Send(events)

	if testHttpServer.Body != `{"api_key":"abc123","expires":1679924153,"signature":"8eTAu1FgPe4568Nxxch8mT+UqpY=","events":[{"stream":"haproxy.requests"}]}` {
		t.Error("Incorrect body:", testHttpServer.Body)
	}
}

func TestSendMultipleEvents(t *testing.T) {
	client := NewTestClient()
	events := make([]Event, 0)

	events = append(events, Event{Stream: "haproxy.requests"})
	events = append(events, Event{Stream: "haproxy.load-time", Weight: 3.4})
	client.Send(events)

	if testHttpServer.Body != `{"api_key":"abc123","expires":1679924153,"signature":"8eTAu1FgPe4568Nxxch8mT+UqpY=","events":[{"stream":"haproxy.requests"},{"stream":"haproxy.load-time","weight":3.4}]}` {
		t.Error("Incorrect body:", testHttpServer.Body)
	}
}

func TestSendCorrectMimeType(t *testing.T) {
	client := NewTestClient()
	events := make([]Event, 0)

	events = append(events, Event{Stream: "haproxy.requests"})
	client.Send(events)

	if testHttpServer.ContentType != "application/json" {
		t.Error("Incorrect Content-Type:", testHttpServer.ContentType)
	}
}

func TestHandleInvalidAuthentication(t *testing.T) {
	// Statsy will return an error string the returned json if an error occurs
	client := NewTestClient()
	events := make([]Event, 0)

	testHttpServer.MockStatusCode = 406
	testHttpServer.MockBody = "{\"error\":\"Invalid signature\"}"
	err := client.Send(events)
	correctError := errors.New("Invalid signature")
	if err.Error() != correctError.Error() {
		t.Error("Incorrect error string received:", err.Error(), "expected", correctError.Error())
	}
}

func TestSendTimestamp(t *testing.T) {
	client := NewTestClient()
	events := make([]Event, 0)

  eventTime := time.Date(2013, 3, 27, 11, 03, 18, 0, time.UTC)
	events = append(events, Event{Stream: "haproxy.requests", Time: eventTime})
	client.Send(events)

	if testHttpServer.Body != `{"api_key":"abc123","expires":1679924153,"signature":"8eTAu1FgPe4568Nxxch8mT+UqpY=","events":[{"stream":"haproxy.requests","timestamp":1364382198}]}` {
		t.Error("Incorrect body:", testHttpServer.Body)
	}
}

func TestIncrement(t *testing.T) {
  client := NewTestClient()

  client.Increment("haproxy.bytes", 634)

  if testHttpServer.Body != `{"api_key":"abc123","expires":1679924153,"signature":"8eTAu1FgPe4568Nxxch8mT+UqpY=","events":[{"stream":"haproxy.bytes","weight":634}]}` {
    t.Error("Incorrect body:", testHttpServer.Body)
  }
}
