package cmpp2

import (
	"testing"
)

func Test_cmpp2(t *testing.T) {
	c, err := NewCmpp2Connect("127.0.0.1", 7890)
	if err != nil {
		t.Fatal(err)
	}
	c.close()

}
