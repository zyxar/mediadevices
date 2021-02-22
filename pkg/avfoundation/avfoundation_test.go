// +build darwin

package avfoundation_test

import (
	"testing"

	"github.com/pion/mediadevices/pkg/avfoundation"
	"github.com/pion/mediadevices/pkg/prop"
)

func TestVideoDevice(t *testing.T) {
	devices, err := avfoundation.Devices(avfoundation.Video)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Skip("no video device found")
	}
	session, err := avfoundation.NewSession(devices[0])
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	var mediaProperty prop.Media
	for _, p := range session.Properties() {
		t.Log(p.FrameFormat, p.Width, p.Height)
		mediaProperty = p
	}
	rc, err := session.Open(mediaProperty)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	_, _, err = rc.Read()
	if err != nil {
		t.Errorf("rc.Read: %v", err)
	}
}
