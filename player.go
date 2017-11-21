// PLAYER OBJECT STRUCT AND METHODS
package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

var commandList = map[string]string{"pause": "p", "up": "\x1b[A", "down": "\x1b[B", "left": "\x1b[D", "right": "\x1b[C"}

// Player ...
type Player struct {
	Playing  bool
	FileName string
	Handler  *exec.Cmd
	PipeIn   io.WriteCloser
}

// NewPlayer ...
func NewPlayer() *Player {
	return &Player{
		Playing: false,
	}
}

// Start ...
func (p *Player) Start(file string) (err error) {
	p.Playing = true
	p.FileName = file
	p.Handler = exec.Command("omxplayer", "-b", "-o", "hdmi", p.FileName)
	p.Handler.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.PipeIn, err = p.Handler.StdinPipe()
	if err == nil {
		p.Handler.Stdout = os.Stdout
		err = p.Handler.Start()
	}
	return
}

// End ...
func (p *Player) End() error {
	pgid, err := syscall.Getpgid(p.Handler.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
		p.FileName = ""
		p.Playing = false
	}
	return err
}

// SendCommand ...
func (p *Player) SendCommand(command string) error {
	_, err := p.PipeIn.Write([]byte(commandList[command]))
	return err
}
