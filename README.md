# Narra Assistant
> A software analytical desktop app

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

- **Frontend**:  
  - [![React](https://img.shields.io/badge/React-%2320232a.svg?logo=react&logoColor=%2361DAFB)](https://react.dev/)  
  - [![TypeScript](https://img.shields.io/badge/TypeScript-%233178C6.svg?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)  

- **Backend**:  
  - [![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?logo=go&logoColor=white)](https://go.dev/)  

- **AI & Pipeline**:  
  - [![Gemini](https://img.shields.io/badge/Gemini%20LLM-%234285F4.svg?logo=google&logoColor=white)](https://deepmind.google/technologies/gemini/)  
  - [![RAG](https://img.shields.io/badge/RAG%20Pipeline-%236DB33F.svg?logo=vectorworks&logoColor=white)](#)


<img width="1917" height="1123" alt="image" src="https://github.com/user-attachments/assets/36b26854-a364-4c5b-bb45-61da3cf2e065" />


<img width="1914" height="1119" alt="Screenshot 2025-09-14 165051" src="https://github.com/user-attachments/assets/9df62443-4133-4cc2-88c0-55127619abfe" />



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
   - Go compiler --> [here](https://go.dev/doc/install)
   - Node.js --> [here](https://nodejs.org/en/download/)

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
