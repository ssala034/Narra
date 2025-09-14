import React, { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { Context } from "../../Context/Context";
import { assets } from "../../assets/assets";
import "./Home.css";

interface ContextType {
  projectPath: string;
  setProjectPath: (path: string) => void;
  LoadEmbeddings: (path: string) => Promise<void>;
}

const Home: React.FC = () => {
  const { projectPath, setProjectPath, LoadEmbeddings } = useContext<ContextType>(Context);
  const navigate = useNavigate();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const path = e.target.value;
    setProjectPath(path); 
    console.log("Project Path Submitted:", path);
  };

  const handleSendClick = async () => {
    if (projectPath.trim()) {
      console.log("Project Path Submitted:", projectPath);
      
      LoadEmbeddings(projectPath);
      
      navigate("/loading");
    }
  };

  return (
    <div className="home-container">
      <h1 className="home-title">Narra Assistant</h1>
      <p className="home-subtitle">
        Welcome to the new way to analyze your local projects.  
        Enter your project repository path to analyze it and learn more about your project.
      </p>
      <div className="home-input-container">
        <input 
          className="home-input" 
          placeholder="Enter project path..." 
          value={projectPath}
          onChange={handleInputChange}
        />
        {projectPath && (
          <img 
            className="home-send-icon" 
            src={assets.send_icon} 
            alt="Send" 
            onClick={handleSendClick}
          />
        )}
      </div>
    </div>
  );
};

export default Home;
