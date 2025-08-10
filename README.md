# Narra
A software analytical desktop app

## About

Wails template which includes: Vite, React, TS, TailwindCSS out of the box.

Build with `Wails CLI v2.0.0`.

To use this [template](https://wails.io/docs/community/templates):
```shell
wails init -n "Your Project Name" -t https://github.com/hotafrika/wails-vite-react-ts-tailwind-template
cd frontend/src
npm install
```

[Here](scripts) you can find useful scripts for building on different platforms and Wails CLI installation.

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this in your
browser and connect to your application.

## Building

To build a redistributable, production mode package, use `wails build`.



**Todo**
Plan is get it to work with current small ones then add Chromadb for the embeddings for larger ones !!!
    - it may wokr with regualr embedding but chroma is more sophisticated
    - the loading may actually be long, may chroma can help?!
since gemini response is slow need to have a nice loading thing that most chats have
Understand the pipeline you made fully
while testing was getting some error
```
2025/08/10 13:58:07 Error creating embedding for document cc9aaaeff37e3394: failed to create embedding: proto: field google.ai.generativelanguage.v1beta.Part.text contains invalid UTF-8
```

Try to fix it but good start, just make sure this current version works, then extrapolate to chroma db (next level)
even just parsing the documents take forever
make sure you only parse `source code` so that it doesn't take for every(2025/08/10 14:31:07 Processing document 90901/264797 too long)

Also the amount of tokens may be to small for a large project
Try to use some optimization techniques so that i don't ahve to actually parse 200000 files
    - might try to reascrh things so I can just judge based on a limited set
    - or figure out the most import set and only use that (but need to know what is the most important set??? optimizaiton)
    - also because i may not be able to send that much data even if vectorizing it to gemini

currently very slow reading everything





_____


AI:
since the AI will neeed to work on large repostiories, I must do a properly vector dataase using chroma db

as soon as they put in the filepath, I should have the pipeline ready to go

Need to make sure that if question is out of context, we say so!!

UI:
make it dark themed
make it responsive when the users changes the size
make it collapsable so they can close if they don't need it and then re-open to same size it was before
might need to make new chats with sql so save info ??
make sure that everything slides when going to next section of page like in slides powerpoint


Only after done on:
try light and dark themed



august 9th:
installed : three.js library(npm install three @react-three/fiber@8 @react-three/drei@9 @react-spring/three)


august 10th

current pipeline works for small doucment <10,000
try to use chroma db so that it works for large scale and cloud base  !!!!
https://github.com/amikos-tech/chroma-go



