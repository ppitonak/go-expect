// Copyright 2018 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package expect

import (
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func TestExpect(t *testing.T) {
	t.Parallel()

	c, err := NewTestConsole(t)
	if err != nil {
		t.Errorf("Expected no error but got'%s'", err)
	}
	defer c.Close()

	cmd := exec.Command("go", "run", "./cmd/prompt/main.go")
	cmd.Stdin = c.Stdin()
	cmd.Stdout = c.Stdout()
	cmd.Stderr = c.Stdout()

	go func() {
		c.Expect("What is 1+1?")
		c.SendLine("2")
		c.Expect("What is Netflix backwards?")
		c.SendLine("xilfteN")
		c.ExpectEOF()
	}()

	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
}

func TestExpectOutput(t *testing.T) {
	t.Parallel()

	c, err := NewTestConsole(t)
	if err != nil {
		t.Errorf("Expected no error but got'%s'", err)
	}
	defer c.Close()

	cmd := exec.Command("go", "run", "./cmd/prompt/main.go")
	cmd.Stdin = c.Stdin()
	cmd.Stdout = c.Stdout()
	cmd.Stderr = c.Stdout()

	go func() {
		c.Expect("What is 1+1?")
		c.SendLine("3")
		c.ExpectEOF()
	}()

	expected := "exit status 1"
	err = cmd.Run()
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s' but got '%s' instead", expected, err)
	}
}

func TestEditor(t *testing.T) {
	c, err := NewConsole()
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	defer c.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}

	cmd := exec.Command("vi", file.Name())
	cmd.Stdin = c.Stdin()
	cmd.Stdout = c.Stdout()
	cmd.Stderr = c.Stdout()

	go func() {
		c.Send("iHello world\x1b")
		c.SendLine(":w")
		c.SendLine(":q!")
		c.ExpectEOF()
	}()

	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}

	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if string(data) != "Hello world\n" {
		t.Errorf("Expected '%s' to equal '%s'", string(data), "Hello world\n")
	}
}

func ExampleConsole() {
	c, err := NewConsole()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command("prompt")
	cmd.Stdin = c.Stdin()
	cmd.Stdout = c.Stdout()
	cmd.Stderr = c.Stdout()

	go func() {
		c.Expect("What is 1+1?")
		c.SendLine("2")
		c.Expect("What is Netflix backwards?")
		c.SendLine("xilfteN")
		c.ExpectEOF()
	}()

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
