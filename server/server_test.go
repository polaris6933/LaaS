package main

import (
	"LaaS/life"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

const test_user = "test_user"
const test_password = "1234"
const test_session = "test_session0"

func assert(expected, actual string, t *testing.T) {
	if expected != actual {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func getTestSession() *session {
	u := newUser(test_user, test_password)
	s := NewSession("test_session", u)
	nl, _ := life.NewLife("pulsar")
	s.currState = nl
	return s
}

func getTestServer() *Server {
	s := NewServer()
	s.Register(test_user, test_password)
	for i := 0; i < 10; i++ {
		s.Add(test_user, fmt.Sprintf("test_session%d", i))
	}
	return s
}

func TestSessionRunProperlyStops(t *testing.T) {
	t.Parallel()
	s := getTestSession()
	s.run()
	time.Sleep(time.Second)
	s.stop()
	time.Sleep(time.Second)
	if s.isRunning {
		t.Fatal("expected session to stop running, it keeps going")
	}
}

func TestSessionAuthorize(t *testing.T) {
	t.Parallel()
	s := getTestSession()
	if !s.authorize(test_user) {
		t.Fatal("authorization with correct credentials failed")
	}
	if s.authorize("adsf") {
		t.Fatal("authorization with incorrect credentials failed")
	}
}

func TestSessionIndex(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	goodIdx := s.sessionIndex("test_session3")
	expectedIdx := 3
	if goodIdx != expectedIdx {
		t.Fatalf("expected index %d, got: %d", expectedIdx, goodIdx)
	}

	badIdx := s.sessionIndex("test_session33")
	expectedIdx = -1
	if badIdx != expectedIdx {
		t.Fatalf("expected index %d, got: %d", expectedIdx, badIdx)
	}
}

func TestStartProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Start(test_user, test_session, "pulsar")
	expected := "successfully started session " + test_session
	assert(result, expected, t)

	time.Sleep(time.Second)
	if !s.sessions[0].isRunning {
		t.Fatal("session is not running")
	}
}

func TestStartBadSession(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	session := "test_sessionX"
	result := s.Start(test_user, session, "pulsar")
	assert(result, noSession(session), t)
}

func TestStartRunnigSession(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	time.Sleep(time.Second)
	result := s.Start(test_user, test_session, "blinker")
	expected := "session " + test_session + " is already running"
	assert(result, expected, t)
}

func TestStartBadUser(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	user := "userX"
	result := s.Start(user, test_session, "pulsar")
	assert(result, notAuthorized(user), t)
}

func TestStartBadConfig(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	config := "asdf"
	result := s.Start(test_user, test_session, config)
	expected := "the configuration you specified does not exist"
	assert(result, expected, t)
}

func TestResumeProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	time.Sleep(time.Second)
	s.Stop(test_user, test_session)
	result := s.Resume(test_user, test_session)
	expected := "successfully resumed session " + test_session
	assert(result, expected, t)
	time.Sleep(time.Second)
	if !s.sessions[0].isRunning {
		t.Fatal("expected session to be running, it isn't")
	}
}

func TestResumeRunningSession(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	time.Sleep(time.Second)
	result := s.Start(test_user, test_session, "blinker")
	expected := "session " + test_session + " is already running"
	assert(result, expected, t)
}

func TestStopProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	time.Sleep(time.Second)
	result := s.Stop(test_user, test_session)
	expected := "session " + test_session + " successfully stopped"
	assert(result, expected, t)
}

func TestStopStopped(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	time.Sleep(time.Second)
	s.Stop(test_user, test_session)
	result := s.Stop(test_user, test_session)
	expected := "session " + test_session + " is already stopped"
	assert(result, expected, t)
}

func TestListProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := len(strings.Split(s.List(), "\n"))
	expectedLines := 12
	assert(strconv.Itoa(result), strconv.Itoa(expectedLines), t)
}
