import React, { useEffect, useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { ThreeCircles } from "react-loader-spinner";
import { Context } from "../../Context/Context";
import "./Loading.css";

interface ContextType {
  loadingData: boolean;
}

const Loading: React.FC = () => {
  const navigate = useNavigate();
  const [currentSubtitle, setCurrentSubtitle] = useState(0);
  const { loadingData } = useContext<ContextType>(Context);

  const subtitles = [
    "Setting up your development environment...",
    "Building pipeline...",
    "Analyzing project structure...",
    "Indexing documents...",
    "Creating embeddings...",
    "Initializing RAG system...",
    "Preparing AI assistant...",
    "Almost ready..."
  ];

  useEffect(() => {
    // Navigate to main only when loadingData is false
    if (!loadingData) {
      const timer = setTimeout(() => {
        navigate("/main");
      }, 1000); // Small delay to show "Almost ready..." message

      return () => clearTimeout(timer);
    }
  }, [loadingData, navigate]);

  useEffect(() => {
    // Change subtitle every 1.5 seconds
    const subtitleTimer = setInterval(() => {
      setCurrentSubtitle((prev) => (prev + 1) % subtitles.length);
    }, 1500);

    return () => clearInterval(subtitleTimer);
  }, [subtitles.length]);

  return (
    <div className="loading-container">
      <h2 className="loading-title">Analyzing Your Project...</h2>
      <ThreeCircles
        visible={true}
        height="120"
        width="120"
        color="#9a2740"
        ariaLabel="three-circles-loading"
        wrapperStyle={{}}
        wrapperClass="loading-spinner"
        outerCircleColor="#9a2740"
        middleCircleColor="#221f20"
        innerCircleColor="#ffffff"
      />
      <p className="loading-subtitle">
        {subtitles[currentSubtitle]}
      </p>
    </div>
  );
};

export default Loading;
