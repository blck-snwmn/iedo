package iedo

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-shellwords"
	"github.com/mattn/go-tty"
	"golang.org/x/xerrors"
)

func getEditorEnv() ([]string, error) {
	p := shellwords.NewParser()
	// EDITOR=vi
	// EDITOR=code -w
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	eParams, err := p.Parse(editor)
	if err != nil {
		return nil, err
	}
	return eParams, nil
}

func openEditor(fileName string) error {
	eParams, err := getEditorEnv()
	if err != nil {
		return err
	}
	eParams = append(eParams, fileName)

	tty, err := tty.Open()
	if err != nil {
		return err
	}
	defer tty.Close()

	cmd := exec.Command(eParams[0], eParams[1:]...)
	cmd.Stdin = tty.Input()
	cmd.Stdout = tty.Output()
	cmd.Stderr = tty.Output()

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func output(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = io.Copy(os.Stdout, f); err != nil {
		return err
	}
	return nil
}

// readInput return (filename, error)
func readInput(in io.Reader) (string, error) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		return "", err
	}
	defer func() {
		f.Close()
	}()
	if _, err := io.Copy(f, in); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

func Run() error {
	in := os.Stdin
	if isatty.IsTerminal(in.Fd()) {
		return xerrors.New("no support for termila input")
	}
	tmp, err := readInput(in)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tmp)
	}()
	if err = openEditor(tmp); err != nil {
		return nil
	}
	if err = output(tmp); err != nil {
		return err
	}
	return nil
}
