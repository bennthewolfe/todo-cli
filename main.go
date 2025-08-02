package main

func main() {
	// This is the main function where the program starts execution.
	// You can add your code here to implement the desired functionality.
	todoList := TodoList{}

	todoList.add("Buy groceries")
	todoList.add("Walk the dog")
	todoList.add("Read a book")

	todoList.delete(1) // Delete the second item

	todoList.view("json")
}
