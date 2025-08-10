package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
	// "google.golang.org/genai" // Issue with the import, might be how the go.mod owrks ???!?!?!
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

/*
Logic for buiding file system, will be used inside the RAG pipeline
Still need to do data indexing and vector store creation
(might do on different files???)
*/



/*
Logic for building file system, will be used inside the RAG pipeline
*/
type FileSystemNode struct {
	Name     string
	Path     string
	Parent   *FileSystemNode
	IsDir    bool
	Children []*FileSystemNode
	Content  []byte // Content is used for files, empty for directories
}

func buildChildren(dirPath string, parent *FileSystemNode) []*FileSystemNode {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Consider logging the error instead of panicking in production
		panic(err)
	}

	var children []*FileSystemNode

	for _, entry := range entries {
		// Build the full path for this entry
		fullPath := filepath.Join(dirPath, entry.Name())

		node := &FileSystemNode{
			Name:   entry.Name(),
			Path:   fullPath,
			Parent: parent,
			IsDir:  entry.IsDir(),
		}

		if entry.IsDir() {
			// Recursively build children for subdirectories
			node.Children = buildChildren(fullPath, node)
		} else {
			// Read file content for regular files
			content, err := os.ReadFile(fullPath)
			if err != nil {
				// Handle error appropriately - could log and continue or panic
				// For now, keeping consistent with your error handling approach
				panic(err)
			}
			node.Content = content
		}

		children = append(children, node)
	}

	return children
}

// Helper function to build the root node and start the recursive process
func BuildFileSystemTree(rootPath string) (*FileSystemNode, error) {
	// Check if the root path exists and is a directory
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, os.ErrInvalid // rootPath is not a directory
	}

	root := &FileSystemNode{
		Name:   filepath.Base(rootPath),
		Path:   rootPath,
		Parent: nil,
		IsDir:  true,
	}

	// Build children recursively
	root.Children = buildChildren(rootPath, root)

	return root, nil
}

// PrintFileSystemTree prints the directory structure in a tree format
func PrintFileSystemTree(node *FileSystemNode) {
	printNode(node, "", true, true)
}

// printNode is a recursive helper function that handles the tree printing logic
func printNode(node *FileSystemNode, prefix string, isLast bool, isRoot bool) {
	if node == nil {
		return
	}

	// Print the current node
	if isRoot {
		// For root node, just print the name
		fmt.Printf("%s\n", node.Name)
	} else {
		// Choose the appropriate tree characters
		var connector string
		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		// Add directory indicator for directories
		name := node.Name
		if node.IsDir {
			name += "/"
		}

		fmt.Printf("%s%s%s\n", prefix, connector, name)
	}

	// If it's a directory, print its children
	if node.IsDir && len(node.Children) > 0 {
		// Determine the prefix for children
		var childPrefix string
		if isRoot {
			childPrefix = ""
		} else if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}

		// Print each child
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1
			printNode(child, childPrefix, isLastChild, false)
		}
	}
}

/*

Connection to google gemini AI API logic
need to test

*/

func (a *App) connectToGeminiAPI() {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	// defer client.Close() // carful with this on when to close the client, it should be closed when you are done with it ONLYS

	if err != nil {
		log.Fatal(err)
	}

	// TRY Disables thinking (decide later base on the tokens I have what to do with the thinking budget)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("What type of doctor is a optologist? What organ do they study?"),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Gemini Response .....")
	fmt.Println(result.Text)
	fmt.Println("End of Gemini Response")
}
