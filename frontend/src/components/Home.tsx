import React from "react";
import styled from "styled-components";

// Full-screen container with wave gradient background
const Container = styled.div`
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: white;
  text-align: center;
  padding: 0 1rem;
  position: relative;
  overflow: hidden;
  background: linear-gradient(
    135deg,
    #090607 0%,
    #221f20 30%,
    #090607 65%,
    #9a2740 100%
  );

  &::before {
    content: "";
    position: absolute;
    bottom: -50px; /* pushes wave down */
    left: 0;
    width: 100%;
    height: 200px; /* wave height */
    background: radial-gradient(
        circle at 70% 40%,
        rgba(154, 39, 64,0),
        transparent 70%
      ),
      radial-gradient(
        circle at 20% 80%,
        rgba(154, 39, 64, 1),
        transparent 80%
      );
    clip-path: path(
      "M0,160 C480,300 960,40 1440,160 L1440,320 L0,320 Z"
    ); /* custom uneven wave */
    opacity: 0.8;
  }
`;

// Narra title
const Title = styled.h1`
  font-family: "Handlee", cursive;
  font-size: 3rem;
  margin-bottom: 0.5rem;
  text-shadow: 0 0 12px rgba(255, 255, 255, 0.6);
  z-index: 1;
`;

// Subtitle styling
const Subtitle = styled.p`
  font-size: 1.1rem;
  max-width: 500px;
  margin-bottom: 1.5rem;
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.85);
  z-index: 1;
`;

// Input box
const Input = styled.input`
  padding: 0.8rem 1rem;
  font-size: 1rem;
  border-radius: 8px;
  border: none;
  outline: none;
  width: 300px;
  text-align: center;
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.2);
  z-index: 1;

  &::placeholder {
    color: rgba(255, 255, 255, 0.6);
  }
`;

const Home: React.FC = () => {
  return (
    <Container>
      <Title>Narra Assistant</Title>
      <Subtitle>
        Welcome to the new way to analyze your local projects.  
        Enter your project repository path to analyze it and learn more about your project.
      </Subtitle>
      <Input placeholder="Enter project path..." />
    </Container>
  );
};

export default Home;
