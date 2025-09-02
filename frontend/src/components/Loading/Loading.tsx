import React from "react";
import "./Loading.css";

interface LoaderProps {
  size?: number;
}
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! Loader is broken right now!!!!!!!!!
const Loading: React.FC<LoaderProps> = ({ size = 80 }) => {
  return (
    <div
      className="loader"
      style={{
        width: size,
        height: size,
      }}
    >
      <div className="loader-circle" />
    </div>
  );
};

export default Loading;
