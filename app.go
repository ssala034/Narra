package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
FileSystemNode represents a hierarchical file system structure used for document indexing.
This implementation supports both directories and files, with content stored for text files
that will be processed by the RAG pipeline for document analysis and vector embeddings.
*/
type FileSystemNode struct {
	Name     string
	Path     string
	Parent   *FileSystemNode
	IsDir    bool
	Children []*FileSystemNode
	Content  []byte // Content stores file data for text files, empty for directories
}

func buildChildren(dirPath string, parent *FileSystemNode) []*FileSystemNode {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Log error and continue in production environments
		panic(err)
	}

	var children []*FileSystemNode

	for _, entry := range entries {
		// Construct the full path for this entry
		fullPath := filepath.Join(dirPath, entry.Name())

		node := &FileSystemNode{
			Name:   entry.Name(),
			Path:   fullPath,
			Parent: parent,
			IsDir:  entry.IsDir(),
		}

		if entry.IsDir() {
			// Recursively process subdirectories
			node.Children = buildChildren(fullPath, node)
		} else {
			// Read and store file content for regular files
			content, err := os.ReadFile(fullPath)
			if err != nil {
				// Handle file read errors - could implement logging instead of panic
				panic(err)
			}
			node.Content = content
		}

		children = append(children, node)
	}

	return children
}

// BuildFileSystemTree constructs a complete file system tree starting from the specified root path.
func BuildFileSystemTree(rootPath string) (*FileSystemNode, error) {
	// Validate that the root path exists and is a directory
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, os.ErrInvalid // Root path must be a directory
	}

	root := &FileSystemNode{
		Name:   filepath.Base(rootPath),
		Path:   rootPath,
		Parent: nil,
		IsDir:  true,
	}

	// Build the complete file system tree recursively
	root.Children = buildChildren(rootPath, root)

	return root, nil
}

// PrintFileSystemTree displays the directory structure in a tree format.
func PrintFileSystemTree(node *FileSystemNode) {
	printNode(node, "", true, true)
}

// printNode recursively renders the tree structure with appropriate formatting.
func printNode(node *FileSystemNode, prefix string, isLast bool, isRoot bool) {
	if node == nil {
		return
	}

	// Render the current node with appropriate formatting
	if isRoot {
		// Display root node without tree connectors
		fmt.Printf("%s\n", node.Name)
	} else {
		// Select appropriate tree connector symbols
		var connector string
		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		// Add visual indicator for directories
		name := node.Name
		if node.IsDir {
			name += "/"
		}

		fmt.Printf("%s%s%s\n", prefix, connector, name)
	}

	// Recursively process child nodes for directories
	if node.IsDir && len(node.Children) > 0 {
		// Calculate prefix for child nodes
		var childPrefix string
		if isRoot {
			childPrefix = ""
		} else if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}

		// Render each child node with proper formatting
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1
			printNode(child, childPrefix, isLastChild, false)
		}
	}
}
