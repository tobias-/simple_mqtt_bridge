package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoPassword(t *testing.T) {
	urlUnderTest := "ssl://foo@localhost/bar"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "foo:bar", username)
	assert.Equal(t, "", password)
	assert.Equal(t, "ssl://localhost:8883", parsedUrl)
}

func TestColonInPassword(t *testing.T) {
	urlUnderTest := "ssl://foo:bar:barbar@localhost/bar"
	_, err := ParseMqttArgs(urlUnderTest, t.Name())
	assert.NotNil(t, err, urlUnderTest, " should not be parseable")
}

func TestNoVhost(t *testing.T) {
	urlUnderTest := "ssl://foo:bar@localhost/"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "ssl://localhost:8883", parsedUrl)
	assert.Equal(t, "foo", username)
	assert.Equal(t, "bar", password)
}

func TestNoTopic(t *testing.T) {
	urlUnderTest := "ssl://foo:bar@localhost/vhost"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "ssl://localhost:8883", parsedUrl)
	assert.Equal(t, "foo:vhost", username)
	assert.Equal(t, "bar", password)
}

func TestNoHost(t *testing.T) {
	urlUnderTest := "ssl://foo:bar@/vhost"
	_, _, _, err := parseUrl(urlUnderTest, t.Name())
	assert.NotNil(t, err, urlUnderTest, " should not be parseable")
}

func TestSunnyDay(t *testing.T) {
	urlUnderTest := "ssl://foo:bar@localhost/vhost"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "ssl://localhost:8883", parsedUrl)
	assert.Equal(t, "foo:vhost", username)
	assert.Equal(t, "bar", password)
}

func TestWss(t *testing.T) {
	urlUnderTest := "wss://foo:bar@localhost/vhost"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "wss://localhost:8883", parsedUrl)
	assert.Equal(t, "foo:vhost", username)
	assert.Equal(t, "bar", password)
}

func TestTcp(t *testing.T) {
	urlUnderTest := "tcp://foo:bar@localhost/vhost"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "tcp://localhost:1883", parsedUrl)
	assert.Equal(t, "foo:vhost", username)
	assert.Equal(t, "bar", password)
}

func TestWs(t *testing.T) {
	urlUnderTest := "ws://foo:bar@localhost/vhost"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "ws://localhost:1883", parsedUrl)
	assert.Equal(t, "foo:vhost", username)
	assert.Equal(t, "bar", password)
}

func TestWithPort(t *testing.T) {
	urlUnderTest := "ssl://foo:bar@hairdresser.cloudmqtt.com:12345/"
	parsedUrl, username, password, err := parseUrl(urlUnderTest, t.Name())
	assert.Nil(t, err, urlUnderTest, " should be parseable")
	assert.Equal(t, "ssl://hairdresser.cloudmqtt.com:12345", parsedUrl)
	assert.Equal(t, "foo", username)
	assert.Equal(t, "bar", password)
}
