package main

func main() {
	// This is the main function where the program starts execution.
	// You can add your code here to implement the desired functionality.
	todoList := TodoList{}

	storage := NewStorage[TodoList]("todos.json")
	loadedList, err := storage.Load()
	if err != nil {
		panic(err)
	}
	todoList = loadedList

	todoList.add("Buy groceries")
	todoList.add("Walk the dog")
	todoList.add("Read a book")

	todoList.delete(1) // Delete the second item

	storage.Save(todoList)

	todoList.view("table")
}
