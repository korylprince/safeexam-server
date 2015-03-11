package api

import (
	"sync"
	"time"
)

// CodeGenerator is an interface for an arbitrary code generation backend.
type CodeGenerator interface {
	// Generate returns the current code from the backend.
	// If the backend malfunctions, code will be empty and error will be non-nil.
	Generate() (code string, expires time.Time, err error)
}

//RandomCodeGenerator represents a CodeGenerator that gets its code using crypto/rand
type RandomCodeGenerator struct {

	//Length specifies how many characters the code should be
	Length int

	current  string
	duration time.Duration
	expires  time.Time
	mu       *sync.Mutex
}

//generator runs in a separate goroutine, updating the code every duration
func generator(r *RandomCodeGenerator) {
	timer := time.NewTimer(0)

	for {
		r.mu.Lock()

		// set  values and start timer
		r.current = randString(r.Length)
		timer.Reset(r.duration)
		r.expires = time.Now().Add(r.duration)

		r.mu.Unlock()

		// wait duration
		<-timer.C
	}
}

//NewRandomCodeGenerator returns a new RandomCodeGenerator with the given code length and duration
func NewRandomCodeGenerator(length int, duration time.Duration) *RandomCodeGenerator {
	r := &RandomCodeGenerator{
		Length:   length,
		current:  "old",
		duration: duration,
		expires:  time.Now(),
		mu:       new(sync.Mutex),
	}
	go generator(r)
	return r
}

// Generate returns the current code. err will always be nil.
func (r *RandomCodeGenerator) Generate() (code string, expires time.Time, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.current, r.expires, nil
}
