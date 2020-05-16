# BRAINSTACK 


```
 brainstack is a program who stock your project and it todos and let you track them.

```


# Installation

```shell
  $ go get 

  $ go build -o brainstack main.go

```

#Usage

```shell
   $ ./brainstack <json file>
   $ init "Read json file Content".
   $ show "Shows a table contains your projects and their todos.
   $ done --project <project name> "For now we only support LIFO which means we remove the last todo of your project" .
   $ add --project <project name> --todo <todos here> add new project and it todos .

```

