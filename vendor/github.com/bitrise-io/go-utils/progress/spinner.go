package progress

import (
	"fmt"
	"io"
	"os"
	"time"
	"unicode/utf8"
)

// Spinner ...
type Spinner struct {
	chars  []string
	delay  time.Duration
	writer io.Writer

	active     bool
	lastOutput string
	stopChan   chan bool
}

// NewSpinner ...
func NewSpinner(chars []string, delay time.Duration, writer io.Writer) Spinner {
	return Spinner{
		chars:  chars,
		delay:  delay,
		writer: writer,

		active:   false,
		stopChan: make(chan bool),
	}
}

// NewDefaultSpinner ...
func NewDefaultSpinner() Spinner {
	chars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	delay := 100 * time.Millisecond
	writer := os.Stdout
	return NewSpinner(chars, delay, writer)
}

func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	for _, c := range []string{"\b", " ", "\b"} {
		for i := 0; i < n; i++ {
			fmt.Fprintf(s.writer, c)
		}
	}
	s.lastOutput = ""
}

// Start ...
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true

	go func() {
		for {
			for i := 0; i < len(s.chars); i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.erase()

					out := s.chars[i]
					fmt.Fprint(s.writer, out)
					s.lastOutput = out

					time.Sleep(s.delay)
				}
			}
		}
	}()
}

// Stop ...
func (s *Spinner) Stop() {
	if s.active {
		s.active = false
		s.erase()
		s.stopChan <- true
	}
}
