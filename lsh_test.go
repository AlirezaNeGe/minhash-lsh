package minhashlsh

import (
	"math/rand"
	"testing"
)

func randomSignature(size int, seed int64) []uint64 {
	r := rand.New(rand.NewSource(seed))
	sig := make([]uint64, size)
	for i := range sig {
		sig[i] = uint64(r.Int63())
	}
	return sig
}

func Test_HashKeyFunc16(t *testing.T) {
	sig := randomSignature(2, 1)
	f := hashKeyFuncGen(2)
	hashKey := f(sig)
	if len(hashKey) != 2*2 {
		t.Fatal(len(hashKey))
	}
}

func Test_HashKeyFunc64(t *testing.T) {
	sig := randomSignature(2, 1)
	f := hashKeyFuncGen(8)
	hashKey := f(sig)
	if len(hashKey) != 8*2 {
		t.Fatal(len(hashKey))
	}
}

func Test_MinhashLSH(t *testing.T) {
	f := NewMinhashLSH16(256, 0.6, 3)
	// sig1 is different from sig2 and sig3
	// sig2 and sig3 are identical
	sig1 := randomSignature(256, 1)
	sig2 := randomSignature(256, 2)
	sig3 := randomSignature(256, 2)

	f.Add("sig1", sig1)
	f.Add("sig2", sig2)
	f.Add("sig3", sig3)

	// No keys are searchable without calling Index()
	if len(f.Query(sig3)) != 0 {
		t.Fatal("No keys should be searchable before calling Index()")
	}

	f.Index()
	for i := range f.hashTables {
		// Hash tables should have size 3
		if len(f.hashTables[i]) != 3 {
			t.Fatal(f.hashTables[i])
		}
	}

	found := 0
	for _, key := range f.Query(sig3) {
		if key.(string) == "sig3" || key.(string) == "sig2" {
			found++
		}
	}
	if found != 2 {
		t.Fatal("unable to retrieve inserted keys")
	}
}

func Test_MinhashLSH2(t *testing.T) {
	minhashLsh := NewMinhashLSH16(256, 0.5, 1)
	seed := int64(1)
	numHash := 256
	mh := NewMinhash(seed, numHash)
	words := []string{"hello", "world", "minhash", "one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine", "ten"}
	for _, word := range words {
		mh.Push([]byte(word))
	}
	sig1 := mh.Signature()
	minhashLsh.Add("s1", sig1)
	minhashLsh.Index()
	k, l := minhashLsh.Params()
	t.Logf("Minhash LSH params: k = %d, l = %d", k, l)

	mh = NewMinhash(seed, numHash)
	words = []string{"one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine", "ten"}
	for _, word := range words {
		mh.Push([]byte(word))
	}
	sig2 := mh.Signature()

	results := minhashLsh.Query(sig2)
	t.Log(results)
	if len(results) < 1 {
		t.Fail()
	}
}

func Test_RemoveKey(t *testing.T) {
	minhashLsh := NewMinhashLSH16(256, 0.5, 1)
	seed := int64(1)
	numHash := 256
	mh := NewMinhash(seed, numHash)
	words := []string{"hello", "world", "minhash", "one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine", "ten"}
	for _, word := range words {
		mh.Push([]byte(word))
	}
	// sig1 and sig2 are identical
	sig1 := mh.Signature()
	minhashLsh.Add("s1", sig1)
	minhashLsh.Index()
	k, l := minhashLsh.Params()
	t.Logf("Minhash LSH params: k = %d, l = %d", k, l)

	mh = NewMinhash(seed, numHash)
	words = []string{"hello", "world", "minhash", "one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine", "ten"}
	for _, word := range words {
		mh.Push([]byte(word))
	}
	sig2 := mh.Signature()
	minhashLsh.Add("s2", sig2)
	minhashLsh.Index()
	k, l = minhashLsh.Params()
	t.Logf("Minhash LSH params: k = %d, l = %d", k, l)

	minhashLsh.Remove("s2")
	minhashLsh.Index()

	mh = NewMinhash(seed, numHash)
	words = []string{"one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine", "ten"}
	for _, word := range words {
		mh.Push([]byte(word))
	}
	sig3 := mh.Signature()

	results := minhashLsh.Query(sig3)
	t.Log(results)
	if len(results) < 1 {
		t.Fail()
	}
}
