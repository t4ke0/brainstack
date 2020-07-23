[![Go Report Card](https://goreportcard.com/badge/github.com/TaKeO90/brainstack)](https://goreportcard.com/report/github.com/TaKeO90/brainstack)


# BRAINSTACK 


```
 brainstack is a program who stock your project and it todos and let you track them.

```


# Installation

```shell
  $ go get 

  $ go build .

```

# Usage

```shell
   * CLI version
   $ ./brainstack -runcmd -file <json file>
   $brainstack init "Read json file Content".
   $brainstack show "Shows a table contains your projects and their todos".
   $brainstack done --project <project name> "For now we only support LIFO which means we remove the last todo of your project" .
   $brainstack add --project <project name> --todo <todos here> "add new project and it todos" .
   
   * TUI version
   $ ./brainstack -runtui -file <json file>

```

# TODO ..... 
* need to document the tui programs with keybinding and other stuff .
