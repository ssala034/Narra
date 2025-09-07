import { createContext, useState } from "react";
import { SimpleQuery } from '../../wailsjs/go/main/RAGPipeline'
import { Greet } from '../../wailsjs/go/main/App'

export const Context = createContext();


const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

const ContextProvider = (props) => {

    const [input, setInput] = useState("")
    const [recentPrompt, setRecentPrompt] = useState("")
    const [previousPrompts, setPreviousPrompts] = useState([])
    const [showResult, setShowResult] = useState(false)
    const [loading, setLoading] = useState(false)
    const [resultData, setResultData] = useState("")

    const delayPara = (index, nextWord) => {
        setTimeout(function () {
            setResultData(prev => prev + nextWord)
            
        }, 75 * index)

    }

    const newChat = () => {
        setLoading(false)
        setShowResult(false)

    }



    const onSent = async (prompt) => {
        /* NOTE */
        // gemini function
        // would be from goLang backend !!! IMPORTANT TO FINISH
        // should be asynchronous function



        //  < input >  is the correct input to use

        setResultData("") // previous response removed
        setLoading(true)
        setShowResult(true)
        let response;
        if(prompt !== undefined){
            console.log("Prompt's type is ", typeof prompt)
            response = await SimpleQuery(prompt)

            setRecentPrompt(prompt)
        }else{
            setPreviousPrompts(prev=>[...prev, input]) // might use the sqlite here, cause I don't want it to re-generate a different reponse
            setRecentPrompt(input)
            //response = await runChat(input)
            console.log("Input's type is ", typeof input)
            response = await SimpleQuery(input)
        }

        let responseArray = response.split("**")
        let newResponse = "";
        for(let i = 0; i < responseArray.length; i++){
            if (i === 0 || i % 2 !== 1){
                newResponse += responseArray[i]
            }
            else{
                newResponse += "<b>" + responseArray[i] + "</b>";
            }
        }
        let newResponse2 = newResponse.split("*").join("</br>")

        let newReponseArray = newResponse2.split(" ");
        for(let i = 0; i < newReponseArray.length; i++){
            const nextWord = newReponseArray[i]
            delayPara(i, nextWord + " ")
        }
        setLoading(false)
        setInput("")


    }

    // onSent("What is react js")


    const contextValue = {
        previousPrompts,
        setPreviousPrompts,
        onSent,
        setRecentPrompt,
        recentPrompt,
        showResult,
        loading,
        resultData,
        input,
        setInput,
        newChat
    }


    return (
        <Context.Provider value={contextValue}>
            {props.children}
        </Context.Provider>
    )
}


export default ContextProvider;