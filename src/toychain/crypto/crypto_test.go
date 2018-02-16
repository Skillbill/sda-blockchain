//
// This is just demo code that happens to compile and run.
//
// This software is just plain non-sense, non-idomatic, bugged,
// unpolished and unfinished.
// It's just my studying notes hacked in a pkg!
//
// This stuff is intended to be used to explain things while looking
// decent and simple on the slides.
//
package crypto

import "testing"

func genkeys() (PrivateKey, PublicKey) {
	priv := NewPrivateKey()
	pub := priv.Public()
	return priv, pub
}

func TestSignVerify(t *testing.T) {
	hash := Hash([]byte("payload"))
	hash2 := Hash([]byte("payload2"))

	priv, pub := genkeys()
	priv2 := genkeys()
	sign := priv.Sign(hash)

	if pub.Verify(hash, sign) == false {
		t.Errorf("exp valid; got invalid")
	}

	if pub.Verify(hash, sign2) == true {
		t.Errorf("exp invalid; got valid")
	}
}
