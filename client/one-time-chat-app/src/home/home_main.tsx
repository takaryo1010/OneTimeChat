import React from 'react';
import CreatePopup from './create_popup.tsx';
import './css/home.css';
const Home: React.FC = () => {


    

    return (
        <div className='home'>
        <h1>ようこそOneTimeChatへ!</h1>
        <p>本サービスは使い捨てのチャットルームを提供します</p>

        <CreatePopup />
        
        </div>
    );
    }


export default Home;