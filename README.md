# Narra Assistant
A software analytical desktop app

## Overview

**Narra** is a software analytical desktop application that leverages **Gemini LLM API** and a **Retrieval-Augmented Generation (RAG) pipeline** to deliver context-aware, personalized interactions. By combining a modern frontend with a Go backend, Narra provides a seamless environment for multi-chat conversations while minimizing hallucinations through intelligent prompt handling.

### Features:
- **Personalized Prompts**: Generate responses tailored to individual conversations.
- **Multiple Chats**: Manage several chat sessions in parallel.
- **Context-Aware Responses**: Maintain continuity across prompts and reduce hallucinations.
- **LLM Integration**: Built with Gemini LLM for advanced natural language capabilities.

## Technologies Used

- **Frontend**: React with TypeScript  
- **Backend**: Go with Wails framework  
- **Gemini LLM**: Provides natural language understanding and generation  
- **RAG Pipeline**: Enhances responses with vector embeddings for context retrieval  

![React](https://img.shields.io/badge/React-%2361DAFB.svg?logo=react&logoColor=black)  
![TypeScript](https://img.shields.io/badge/TypeScript-%233178C6.svg?logo=typescript&logoColor=white)  
![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?logo=go&logoColor=white)  

(pictures)


### RAG Pipeline Architecture
```
+------------------+    +--------------+    +------------+    +---------------+    +-------------------------+
|   React frontend | -> |    Wails     | -> |    RAG     | -> |   Gemini API  | -> |   Context-Aware Reponse |
|    (User input)  |    | (Go backend) |    | (Embedding)|    |  (LLM prompt) |    |    (Back to frontend)   |
+------------------+    +--------------+    +------------+    +---------------+    +-------------------------+
```


## Getting Started 🚀

1. Clone the repository
2. Install Wails CLI if not already installed:

   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

   Please also install:
   - Go compiler
   - Node.js 

3. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

4. Setup API Key
    In `rag.go` enter your Gemini API key where it says `Your-API-Key-Here`. If you'd like, you can use a different LLM provider.
    In that case, make sure to set-up LLM connection properly
    
    ```
    func (r *RAGPipeline) StartEmbeddings(rootPath string) []Document 
        geminiAPIKey := "Your-API-Key-Here"
    ```

5. Build the application:

   **Windows**:
   ```bash
   wails build
   ```

6. Usage
    Its important that when entering your project path to use an abosulte path. After it set's up the pipeline, you will good to go. 


## About Wails Project

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this in your
browser and connect to your application.

## Building

To build a redistributable, production mode package, use `wails build`.