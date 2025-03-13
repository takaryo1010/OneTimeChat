import React from 'react';
import { BrowserRouter, Route, Routes, Link } from 'react-router-dom';
import Debug from './debug.tsx';
import Home from './home/home_main.tsx'
import Chat from './chat/chat_main.tsx'

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/chat" element={<Chat/>} />
        <Route path="/" element={<Home/>} />
        <Route path="/debug" element={<Debug />} />
      </Routes>
    </BrowserRouter>      
  );
};

export default App;
