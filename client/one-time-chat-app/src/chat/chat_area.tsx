import React, { useState } from 'react';
import { Box, TextField, IconButton, Typography, Button } from '@mui/material';
import { Send } from '@mui/icons-material';

interface Message {
  sender: string;
  content: string;
  isMe: boolean;
}

interface ChatAreaProps {
  message: Message[];
  sendMessage: (content: string) => void;
}

const ChatArea: React.FC<ChatAreaProps> = ({ message, sendMessage }) => {
  const [inputMessage, setInputMessage] = useState<string>('');
  const [expandedMessages, setExpandedMessages] = useState<boolean[]>(message.map(() => false));

  const changeInputMessage = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputMessage(e.target.value);
  };

  const handleSendMessage = () => {
    if (inputMessage.trim() !== '') {
      sendMessage(inputMessage);
      setInputMessage('');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const toggleExpandMessage = (index: number) => {
    setExpandedMessages((prev) => {
      const newState = [...prev];
      newState[index] = !newState[index];
      return newState;
    });
  };

  const getMessageContent = (content: string, index: number) => {
    const maxLines = 5;
    const maxChars = 200;
    const lines = content.split('\n');
    const charCount = content.length;

    const isLongMessage = lines.length > maxLines || charCount > maxChars;
    const shouldTruncate = isLongMessage && !expandedMessages[index];

    return (
      <Box>
        <Typography
          variant="body1"
          sx={{
            whiteSpace: 'pre-wrap',
            overflow: 'hidden',
            maxHeight: shouldTruncate ? '5em' : 'none',
          }}
        >
          {content}
        </Typography>
        {isLongMessage && (
          <Button size="small" onClick={() => toggleExpandMessage(index)}>
            {expandedMessages[index] ? '閉じる' : '続きを見る'}
          </Button>
        )}
      </Box>
    );
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1, height: '100%' }}>
      {/* メッセージ表示エリア（スクロール可能） */}
      <Box sx={{ flex: 1, overflowY: 'auto', padding: 2, '&::-webkit-scrollbar': { display: 'none' } }}>
  {message.map((m, index) => (
        <Box
          key={index}
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignSelf: m.isMe ? 'flex-end' : 'flex-start', // 自分のメッセージは右寄せ、他人のは左寄せ
            backgroundColor: m.isMe ? '#dcf8c6' : '#ffffff', // 自分は薄緑、他人は白
            color: 'black',
            padding: 1.5,
            borderRadius: 2,
            marginBottom: 1,
            maxWidth: '70%',
            boxShadow: 1,
            marginLeft: m.isMe ? 'auto' : undefined, // 自分のメッセージは右寄せ
            marginRight: !m.isMe ? 'auto' : undefined, // 他の人のメッセージは左寄せ
          }}
        >
          <Typography variant="caption" sx={{ fontWeight: 'bold', color: '#555' }}>
            {m.sender}
          </Typography>
          {getMessageContent(m.content, index)}
        </Box>
      ))}
    </Box>


      {/* メッセージ入力エリア（固定） */}
      <Box sx={{ display: 'flex', alignItems: 'center', borderTop: '1px solid #ccc', padding: 1 ,overflow:`hidden`}}>
        <TextField
          fullWidth
          multiline
          placeholder="メッセージを入力..."
          value={inputMessage}
          onChange={changeInputMessage}
          onKeyDown={handleKeyDown}
          variant="outlined"
          sx={{
            '& .MuiInputBase-root': { borderRadius: '20px', paddingLeft: 2 },
          }}
        />
        <IconButton color="primary" onClick={handleSendMessage} sx={{ marginLeft: 1 }}>
          <Send />
        </IconButton>
      </Box>
    </Box>
  );
};

export default ChatArea;
