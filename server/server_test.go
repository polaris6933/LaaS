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

func assert(actual, expected string, t *testing.T) {
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

func TestRegisterProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	username := "test_user1"
	result := s.Register(username, "asdf")
	expected := "registered user " + username
	assert(result, expected, t)
}

func TestRegisterUsedName(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Register(test_user, "asdf")
	expected := "user " + test_user + " already exists"
	assert(result, expected, t)
}

func TestLoginProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Login(test_user, test_password)
	expected := "user " + test_user + " logged in"
	assert(result, expected, t)
}

func TestLoginNoUser(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	username := "test_userX"
	result := s.Login(username, test_password)
	expected := "user " + username + " does not exist"
	assert(result, expected, t)
}

func TestLoginBadPassword(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Login(test_user, "....")
	expected := "invalid password for " + test_user
	assert(result, expected, t)
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

func TestAddProper(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	session := "new_session"
	result := s.Add(test_user, session)
	expected := "successfully created session " + session
	assert(result, expected, t)
}

func TestAddExisting(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Add(test_user, test_session)
	expected := "session with the name " + test_session + " already exists"
	assert(result, expected, t)
}

func TestKillStopped(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	result := s.Kill(test_user, test_session)
	expected := "session " + test_session + " successfully killed"
	assert(result, expected, t)
}

func TestKillRunning(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	s.Start(test_user, test_session, "pulsar")
	result := s.Kill(test_user, test_session)
	expected := "session " + test_session + " successfully killed"
	assert(result, expected, t)
}

func TestKillNoSession(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	session := "no_session"
	result := s.Kill(test_user, session)
	expected := noSession(session)
	assert(result, expected, t)
}

func TestKillNotAuthorized(t *testing.T) {
	t.Parallel()
	s := getTestServer()
	user := "no_user"
	result := s.Kill(user, test_session)
	expected := notAuthorized(user)
	assert(result, expected, t)
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
