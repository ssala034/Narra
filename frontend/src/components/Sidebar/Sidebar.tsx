import React, { useContext, useState } from "react";
import "./Sidebar.css";
import { assets } from '../../assets/assets';
import { Context } from  '../../Context/Context'


interface ContextType {
  onSent: (prompt: string) => void | Promise<void>;
  previousPrompts: string[]; 
  setRecentPrompt: (prompt: string) => void;
  newChat: any
}


const Sidebar = () => {
  const [extended, setExtended] = useState(false);
  const {onSent, previousPrompts, setRecentPrompt, newChat} = useContext<ContextType>(Context)

  const loadPrompt = async (prompt: any) => {
    setRecentPrompt(prompt)
    await onSent(prompt)
  }


  // Note the previous chats will go away so try to set up a sqlite so you can save previous chats!!!!!


  return (
    <div className="sidebar">
      <div className="top">
        <img onClick={() => setExtended(prev=>!prev)} className="menu" src={assets.menu_icon} alt="" />
        <div onClick={()=> newChat()} className="new-chat">
          <img src={assets.plus_icon} alt="" />
          {extended ? <p>New Chat</p> : null}
        </div>
        {extended ? (
          <div className="recent">
            <p className="recent-title">Recent</p>
            {previousPrompts.map((item, index)=> {
              return (
                <div key={index} onClick={()=>loadPrompt(item)} className="recent-entry">
                <img src={assets.message_icon} alt="" />
                <p>{item.slice(0,18)} ...</p>
              </div>
              )
            })}
            
          </div>
        ) : null}
      </div>
      <div className="bottom">
        <div className="bottom-item recent-entry">
          <img src={assets.question_icon} alt="" />
          {extended ? <p>Help</p> : null}
        </div>
        <div className="bottom-item recent-entry">
          <img src={assets.history_icon} alt="" />
          {extended ? <p>Activities</p> : null}
        </div>
        <div className="bottom-item recent-entry">
          <img src={assets.setting_icon} alt="" />
          {extended ? <p>Settings</p> : null}
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
