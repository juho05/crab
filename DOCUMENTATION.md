# Documentation

The following instructions assume that you have a working [installation](https://github.com/Bananenpro/crab/blob/main/README.md#Installation)
of the _crab_ interpreter.

## Introduction

_crab_ is an interpreted toy programming language, which makes it quiet slow but still usable for simple tasks.

It has a C-like syntax but a dynamic type system, which allows variables to hold values of different types at different times.
Despite its dynamic nature _crab_ tries to catch errors like undefined identifiers, wrong argument or return value counts or `break` statements
outside of loops at parse time.

The standard library currently includes functions for user in- and output, basic handling of lists (e.g. appending, removing items, …) and
simple file operations like reading/writing files or listing all files in a directory.

## Hello World

To beginn writing code in _crab_ you first need to create a new text file. 
The extension doesn't really matter, but the convention is to give your source code file the extension `.cb`.

For a simple 'Hello World' program write the following code in a file called `helloworld.cb`:
```go
func main() {
  println("Hello, World!");
}
```

To execute your code simply call the `crab` interpreter and give it your source file as an argument:
```sh
crab helloworld.cb
```

You should see the following output:
```
Hello, World!
```

### Explanation

The entry point to all _crab_ programs is the `main()` function. It never returns anything and takes no arguments.
The only thing that can be different for its signature is the optional [throws](https://github.com/Bananenpro/crab/blob/main/DOCUMENTATION.md#Exceptions) keyword.

The `main()` function calls the builtin [`println()`](https://github.com/Bananenpro/crab/blob/main/DOCUMENTATION.md#Output) function, an alternative to the [`print()`](https://github.com/Bananenpro/crab/blob/main/DOCUMENTATION.md#Output) function, 
which prints its arguments to `stdout` and appends a newline character.

In general all statements in _crab_ have to end with a `;` similar to other C-like languages.

## Variables

To define a variable in _crab_ simply use the `var` keyword:

```go
var x = 5;
var y; // equivalent to 'var y = null;'
```

Omit the `var` keyword when you want to assign a value to an existing variable:

```go
var x = 5;
// because 'crab' has dynamic typing, the following works:
x = "Hello World!";
```

_crab_ has global and local scope and supports variable shadowing:

```go
// variable x in global scope
var x = 5;

func main() {
	println(x); // 5

	// variable x in local scope
	var x = 7;
	println(x); // 7

	// nested scope
	{
		x = 9;
		var x = 10;
		println(x); // 10
	}

	println(x); // 9
}
```

## Type conversion

Sometimes you need to convert between types, for example when you want to receive numbers from the user.
There are 3 functions for exactly this purpose: `toString()`, `toNumber()` and `toBoolean()`.

### toString()

`toString()` takes a value and converts it to a string. This action will always succeed.

Example:
```go
var x = 5;
x = toString(x); // x is now a string
```

### toNumber()

`toNumber()` takes a string value and tries to convert it to a number. If this action doesn't succeed, `toNumber()` will throw an exception.

`toNumber()` accepts string representations of integers and floating point numbers.

Example:

```go
var x = "5";
try {
	x = toNumber(x); // x is now a number
} catch {
	// couldn't convert x to a number
}
```

### toBoolean()

`toBoolean()` takes a string value and tries to convert it to a boolean. If this action doesn't succeed, `toBoolean()` will throw an exception.

`toBoolean()` accepts "1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false" and "False".

Example:
```go
var x = "true";
try {
	x = toBoolean(x); // x is now a boolean
} catch {
	// couldn't convert x to a boolean
}
```

## Control flow

_crab_ supports the 3 most common control flow constructs:

### If-statement

```javascript
if (condition) {
	// do something
} else if (condition2) {
	// do something else
} else {
	// do something completely different
}
```

### While-loop

```javascript
while (condition) {
	// do something as long as 'condition' is true
}
```

### For-loop

```go
for (var i = 0; i < 10; i++) {
	println(i);
}
```

### Break and continue

_crab_ the `break` and `continue` statements in loops.

- `break`: exit the loop
- `continue`: skip the rest of the current iteration and start over again

## Operators

| symbol | name                | description
|--------|---------------------|--------------------------
| ()     | parentheses         | grouping or function call
| []     | brackets            | list subscript
| ++     | increment           | increments the variable by 1
| --     | decrement           | decrements the variable by 1
| !      | bang                | negates the logical value after it
| -      | unary minus         | multiplies the value after it with -1
| **     | exponentiation      | returns the result of raising the first operand to the power of the second operand
| +      | addition            | adds two values together
| -      | subtraction         | subtracts two values from another
| *      | multiplication      | multiplies two values with each other
| /      | division            | divides a value by another value
| %		 | modulus             | take the modulus of two values
| <      | less                | returns true if the left operand is less than the right one
| >      | greater             | returns true if the left operand is greater than the right one
| <=     | less or equal       | returns true if the left operand is less than or equal to the right one
| >=     | greater or equal    | returns true if the left operand is greater than or equal to the right one
| ==     | equal               | returns true if both operands are equal
| !=     | not equal           | returns true if both operands are not equal
| &&     | logical AND         | returns true if both operands are true
| ||     | logical OR          | returns true if at least one of the operands are true
| ^^     | logical XOR         | returns true if exactly one of the operands is true
| ?:     | ternary conditional | returns either the left or right result depending on the condition
| =      | assignment          | assigns the right value to the left operand
| +=     | assignment          | adds and assigns the right value to the left operand
| -=     | assignment          | subtracts and assigns the right value from the left operand
| *=     | assignment          | multiplies and assigns the right value with the left operand
| /=     | assignment          | divides and assigns the right value from the left operand
| %=     | assignment          | takes the modulus and assigns the right value to the left operand


## Functions

In its most simple form a function can be declared and called as follows:

```go
func test() {
	println("test");
}

func main() {
	test();
}

```

If you want to return values from your function, you will need to tell _crab_ about the number
of values you want to return:

```go
// 0 return values
func test() {
	return;
}

// 1 return value
func test2() 1 {
	return "some string";
}

// 2 return values
func test3() 2 {
	return "some string", 42;
}

func main() {
	test();

	var a = test2(); // a = "some string"

	var b, c = test3(); // b = "some string", c = 42
}
```

_crab_ supports at most **4** return values.

Functions can take arguments by specifying a parameter list between the parantheses:

```go
func add(a, b) 1 {
	return a + b;
}

func main() {
	println(add(1, 2)); // 3
}
```

Functions can be declared inside of other functions with closure support:

```go
func returnAFunction(text) 1 {
	func function() {
		println(text);
	}
	return function;
}

func main() {
	var fn = returnAFunction("Hello, World!");
	fn(); // Hello, World!
}
```

## Strings and lists

### Lists

Lists in _crab_ can hold values of different types and can be dynamically resized:

```javascript
var emptyList = [];
var list = [1, 2, "Hello, World!", [true, false]];

println(list[3]); // Hello, World!

append(list, "new value");
println(list); // [1,2,Hello, World!,[true,false],new value]

remove(list, 3);
println(list); // [1,2,[true,false],new value]
println(len(list)); // 4

concat(list, [3, 4, 5]);
println(list); // [1,2,[true,false],new value,3,4,5]
```

### Strings

You can work with strings similarly as with lists. The only difference is assignment.
Due to the immutable nature of strings, you cannot assign a new character at a specified index.

```go
var helloworld = "Hello, World!";
println(helloworld[4]); // o
helloworld[4] = "y"; // error!
```

### Utility functions

There are a number of builtin utility functions that make working with lists and strings a lot easier:

```go
// will convert all arguments to strings prior to processing
toLower("HelLo, WoRLd!"); // "hello, world!"
toUpper("HelLo, WoRLd!"); // "HELLO, WORLD!"
trim("  \t  Hello, World!       \t  "); "Hello, World!"

// work with both strings and lists
contains("Hello, World!", "!"); // true
contains([1,2,4,5], [2]); // true
indexOf("Hello, World!", "W"); // 7
indexOf([1,2,3,4], [5]); // -1
replace("Hello, World!", "l", "n"); // "Henno, Wornd!"
replace([1,2,3,4], 2, 1); // [1,1,3,4]
split("Hello, World!", ","); // ["Hello", " World!"]
split([1,2,3,4,5], 3); // [[1,2],[4,5]]

// only work with lists
join([1,"hello",false], "-"); // "1-hello-false"
```

## Exceptions

Error handling in _crab_ is done through exceptions.

An exception can be thrown using the `throw` keyword:

```go
func testfunction() throws {
	throw "some value";
}
```

**IMPORTANT:** All exceptions must be explicitly handled by either adding `throws` to the function signature (after the return value count)
or by wrapping the problematic code in a `try...catch` block.

### try...catch

```go
func testfunction() throws {
	throw "some value";
}

func main() {
	try {
		testfunction();
	} catch (e) { // the '(e)' can be omitted, if you don't need the value of the exception
		println(e); // some value
	}
}
```

## User input/output

### Output

There are two functions for printing text to `stdout`: `print()` and `println()`.
They only differ in the way that `println()` appends a newline character while `print()` does not.

You can supply multiple arguments to the `print` functions in order to print them separated by a space.

Example:

```go
println("Hello", "World", 5, true); // Hello World 5 true
```

### Input

You can request input from the user via `stdin` with the `input()` function:

```go
var name = input("What is your name? "); // input always returns a string
```

## File operations

Most file operations can throw exceptions.

### Check if a file exists

Returns true if the specified file exists.

```go
if (fileExist("filepath")) {
	// exists!
}
```

### Read

Returns the full content of the specified file.

```go
var content = readFileText("filepath");
```

### Write

Creates parent directories, if they do not already exists and overrides the file with the provided content.

```go
writeFileText("filepath", "content");
```

### Append to an existing file

Appends the provided content to the specified file.

```go
appendFileText("filepath", "content");
```

### Delete

Deletes the specified file/empty directory.

```go
deleteFile("filepath");
```

### List all files in a directory

Returns a list of all file names in the specified directory.

```go
var files = listFiles("directory");
```

## Other useful functions

### Random number

Generate a random floating point number between _a_ (inclusive) and _b_ (exclusive):

```go
var a = 0;
var b = 100;
var num = random(a, b);
```

Generate a random integer between _a_ (inclusive) and _b_ (exclusive):

```go
var a = 0;
var b = 100;
var num = randomInt(a, b);
```
