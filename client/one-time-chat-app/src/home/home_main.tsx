import React from 'react';
import CreatePopup from './create_popup.tsx';
import JoinPopup from './join_popup.tsx';
import './css/home.css';
import { Box, Button, Typography } from "@mui/material";
import Icon from '@mdi/react';
import { mdiHomePlus, mdiLogin } from '@mdi/js';

const Home: React.FC = () => {
    return (
        <div className="home">
            <h1>ようこそ OneTimeChat へ!</h1>
            <p>本サービスは使い捨てのチャットルームを提供します</p>

            {/* ボタンを横並びに配置 */}
            <Box className="button-container">

            {/* ルーム作成のポップアップ */}
            <CreatePopup />
                {/* ルーム作成ボタン */}
                

                {/* ルーム参加ボタン */}
            <JoinPopup/>
            </Box>
        </div>
    );
}

export default Home;
