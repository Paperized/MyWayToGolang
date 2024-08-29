# MyWayToGolang
Here I will record my small exercises while learning Go. Why not making it public?
Below I will list everything I do and what I learned/was cool to me.

## 1. SumNumbersWithCalcHistory

Small calculator that sums up 2 numbers and register them in a history that can be requested at anytime or at the end of the program.

- Throwback to C/C++ with little pointer operations (Scanf, custom parser that return an integer or a string ptr error)
- After reading some theory of slices and how they work under the hood I implemented a history of my sums operations
- Goofed around the go libraries, I was able to convert string to Numbers, trim strings and add an element at index 0 with the slices package util!
- for loops ranges, for based on a boolean (a while but not a while lol) and some other basic things 

## 2. CatchTheBall

Small game in which there are 5 cups and a ball inside, two players can challenger each others, additionally there is player vs bot, and funny enough bot vs bol lol

- Learned a concise way to do input validation in loop until conditions are satisfied with a for loop
> ```for readStdinTrim(&input); !isValidChoice(input); readStdinTrim(&input) {}```
- Used for the first time structs data type, nested structs. Kinda simple
- Associated some "methods" or rather functions to a struct, they are called Pointer receivers and are similar to how Extensions are used in C#. You can associate functions to a specific type of struct/*struct.
- Having a reference to a function is super easy, just create a variable with func(any,type,in,input)output. I used those to make a game based on events separating input/console prints from game logic.

