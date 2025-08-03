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

	cmdFlags := NewCmdFlag()

	if err := cmdFlags.Execute(&todoList); err != nil {
		panic(err)
	}

	storage.Save(todoList)
}
