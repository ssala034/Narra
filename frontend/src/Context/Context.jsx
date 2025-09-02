import { createContext, useState } from "react";

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
        let promptValue;
        if(prompt !== undefined){
            // response = await runChat(prompt)
            promptValue = true
            setRecentPrompt(prompt)
        }else{
            setPreviousPrompts(prev=>[...prev, input]) // might use the sqlite here, cause I don't want it to re-generate a different reponse
            setRecentPrompt(input)
            //response = await runChat(input)
            promptValue = false;
        }

        // const response = await runChat(input) // return response text
        const responseToRemove = `${promptValue}**React Native** is an open-source framework developed by **Meta** (formerly Facebook) that allows developers to build mobile applications for **iOS** and Android using JavaScript and React. Unlike traditional web apps, React Native apps use native components instead of web views, which means they offer performance close to native apps while sharing a single codebase across platforms. This makes development faster and more efficient, especially for teams building apps for multiple mobile platforms.`
        console.log("Testing response: ", responseToRemove)

        await sleep(2000)


        let responseArray = responseToRemove.split("**")
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