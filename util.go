package cpm

import (
	"fmt"
	"time"
)

type Spinner struct {
	chars  []rune
	delay  time.Duration
	active bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		chars:  []rune{'-', '\\', '|', '/'},
		delay:  100 * time.Millisecond,
		active: true,
	}
}

func (s *Spinner) Start(msg string) {
	fmt.Print(msg)
	s.active = true
	go func() {
		for s.active {
			for _, char := range s.chars {
				fmt.Printf("\b%c", char)
				time.Sleep(s.delay)
			}
		}
	}()
}

func (s *Spinner) Success() {
	s.active = false
	fmt.Print("\b")
	fmt.Print("\033[32m\u2714\033[0m\n")
}

func (s *Spinner) Failure() {
	s.active = false
	fmt.Print("\b")
	fmt.Print("\033[31m\u2717\033[0m\n")
}

func (s *Spinner) Stop() {
	s.active = false
	fmt.Print("\b")
}
