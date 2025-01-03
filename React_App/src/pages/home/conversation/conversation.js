import React, { useState, useEffect, useRef } from 'react';
import {
  eventEmitter,
  sendWebSocketMessage,
} from '../../../services/socket-service';
import { getConversationBetweenUsers } from '../../../services/api-service';

import './conversation.css';

const alignMessages = (userDetails, toUserID) => {
  const { email } = userDetails;
  return email !== toUserID;
};

const scrollMessageContainer = (messageContainer) => {
  if (messageContainer.current !== null) {
    try {
      setTimeout(() => {
        messageContainer.current.scrollTop = messageContainer.current.scrollHeight;
      }, 100);
    } catch (error) {
      console.warn(error);
    }
  }
};

const getMessageUI = (messageContainer, userDetails, conversations) => {
  return (
    <ul ref={messageContainer} className="message-thread-container">
      {conversations.map((conversation, index) => (
        <li
          className={`message ${
            alignMessages(userDetails, conversation.toUserID) ? 'align-right' : ''
          }`}
          key={index}
        >
          {conversation.message}
        </li>
      ))}
    </ul>
  );
};

const getInitiateConversationUI = (userDetails) => {
  if (userDetails !== null) {
    return (
      <div className="message-thread-container start-chatting-banner">
        <p className="heading">
          You haven&apos;t chatted with {userDetails.username} in a while,
          <span className="sub-heading"> Say Hi.</span>
        </p>
      </div>
    );
  }
};

function Conversation(props) {
  const selectedUser = props.selectedUser;
  const userDetails = props.userDetails;

  const messageContainer = useRef(null);
  const [conversation, updateConversation] = useState([]);
  const [messageLoading, updateMessageLoading] = useState(true);

  useEffect(() => {
    if (userDetails && selectedUser) {
      (async () => {
        const conversationsResponse = await getConversationBetweenUsers(userDetails.email, selectedUser.email);

        updateMessageLoading(false);
        if (conversationsResponse.data) {
          updateConversation(conversationsResponse.data);
        } else if (conversationsResponse.response === null) {
          updateConversation([]);
        }
      })();
    }
  }, [userDetails, selectedUser]);

  useEffect(() => {
    const newMessageSubscription = (messagePayload) => {
      if (
        selectedUser !== null &&
        selectedUser.email === messagePayload.fromUserID
      ) {
        updateConversation([...conversation, messagePayload]);
        scrollMessageContainer(messageContainer);
      }
    };

    eventEmitter.on('message-response', newMessageSubscription);

    return () => {
      eventEmitter.removeListener('message-response', newMessageSubscription);
    };
  }, [conversation, selectedUser]);

  const sendMessage = (event) => {
    if (event.key === 'Enter') {
      const message = event.target.value;

      if (message === '' || message === undefined || message === null) {
        alert(`Message can't be empty.`);
      } else if (userDetails.userID === '') {
        this.router.navigate(['/']);
      } else if (selectedUser === undefined) {
        alert(`Select a user to chat.`);
      } else {
        event.target.value = '';

        const messagePayload = {
          fromUserID: userDetails.email,
          message: message.trim(),
          toUserID: selectedUser.email,
        };

        sendWebSocketMessage(messagePayload);
        updateConversation([...conversation, messagePayload]);

        scrollMessageContainer(messageContainer);
      }
    }
  };

 


  const initiateVideoCall = () => {

    window.open(`https://video-chat-1-k9c7.onrender.com/create`, "_blank");
    const messagePayload = {
        fromUserID: userDetails.email,
        message: `Initiating a video call.....................................
        Open this link to connect with me: 
        ..................."https://video-chat-1-k9c7.onrender.com"...........
        ......................................................................
        Wait for 1 minute, I will send you the code to connect with me.`,
        toUserID: selectedUser.email,
    };

    sendWebSocketMessage(messagePayload); 
    updateConversation([...conversation, messagePayload]); 

    scrollMessageContainer(messageContainer); 

    // Redirect to video call page
    // window.location.href = `https://video-chat-1-k9c7.onrender.com`;
   
  };

  if (messageLoading) {
    return (
      <div className="message-overlay">
        <h3>
          {selectedUser !== null && selectedUser.username
            ? 'Loading Messages'
            : ' Select a User to chat.'}
        </h3>
      </div>
    );
  }

  return (
    <div className="app__conversion-container">
      {conversation.length > 0
        ? getMessageUI(messageContainer, userDetails, conversation)
        : getInitiateConversationUI(selectedUser)}

      <div className="app__text-container">
      
        <textarea
          placeholder={`${
            selectedUser !== null ? '' : 'Select a user and'
          } Type your message here`}
          className="text-type"
          onKeyPress={sendMessage}
        ></textarea>

       
       
      </div>
      <div><button
          className="video-call-button"
          onClick={initiateVideoCall}
          title={`Call ${selectedUser ? selectedUser.username : 'user'}`}
          disabled={!selectedUser}
        >
          📞 Video
        </button></div>
      
    </div>
  );
}

export default Conversation;
