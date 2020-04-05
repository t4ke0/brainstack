# BRAINSTACK 
```
  Program that stocks Your Idea as Stack of Ideas.
```


## Usage 

```
  $ go build main.go -o brainstack
  $ ./brainstack -f <CSV FILE HERE> -i <ideas HERE> -add <true if you wanna add these ideas>
  $ ./brainstack -done 

```
## Explanation 


```
  When Adding New ideas or thaughts to the program you each row in the csv file should have the same length ,
  because when you use csv you should respect the number of fields of each line

  $ ./brainstack -f <csv filename> -i <ideas or thaugths here should be separated by comma> -add 
  so number of the ideas separated by comma should be the same in the next addition of the ideas if you wont put the same number of ideas , you can use comma as an additional element per example : 
  $ ./brainstack -f file.csv -i "idea1,idea2,idea3" -add 
  $ ./brainstack -f file.csv -i "idea4,idea5,." -add 
  so you see after idea5 we put comma there to give that row the same number of element of the first row .
```

