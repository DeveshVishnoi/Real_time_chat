import React, {useState} from 'react';
import { withRouter } from 'react-router-dom';

import { isEmailAvailableHTTPRequest, registerHTTPRequest } from "./../../../services/api-service";
import { setItemInLS } from "./../../../services/storage-service";

import './registration.css'

function Registration(props) {
  const [registrationErrorMessage, setErrorMessage] = useState(null);
  const [username, updateUsername] = useState(null);
  const [password, updatePassword] = useState(null);
  const [email, updateEmail] = useState(null)

  let typingTimer = null;


  const handlePasswordChange = async (event) => {
    updatePassword(event.target.value);
  }

  const handleEmailChange = async (event) => {
    updateEmail(event.target.value);
  }

  const handleKeyDownChange = (event) => {
    clearTimeout(typingTimer);
  }

  const handleUsernameChange= async (event)=>{
    updateUsername(event.target.value)
  }
  const handleKeyUpChange = (event) => {
    const email = event.target.value;
    typingTimer = setTimeout( () => {
      checkIfEmailAvailable(email);
    }, 1200);
  }

  const checkIfEmailAvailable = async (username) => {  
    props.displayPageLoader(true);
    const isEmailAvailableResponse = await isEmailAvailableHTTPRequest(username);
    
    props.displayPageLoader(false);
    if (!isEmailAvailableResponse.data) {
      setErrorMessage(isEmailAvailableResponse.message);
    } else {
      setErrorMessage(isEmailAvailableResponse.message);
    }
    updateEmail(username);
  }

  const registerUser = async () => {
    props.displayPageLoader(true);
    const userDetails = await registerHTTPRequest(username, email,password);
    props.displayPageLoader(false);

    
    if (userDetails.statusCode === 201) {
      setItemInLS('userDetails', userDetails.data)
      props.history.push(`/home`)
    } else {
      setErrorMessage(userDetails.message);
    }
  };

  return (
    <div className="app__register-container">
       <div className="app__form-row">
        <label>Username:</label>
        <input type="email" className="email" onKeyDown={handleUsernameChange} />
      </div>

      <div className="app__form-row">
        <label>Email:</label>
        <input type="email" className="email" onChange= {handleEmailChange} onKeyDown={handleKeyDownChange}  onKeyUp={handleKeyUpChange}/>
      </div>
      <div className="app__form-row">
        <label>Password:</label>
        <input type="password" className="password" onChange={handlePasswordChange}/>
      </div>
      <div className="app__form-row">
        <span className="error-message">{registrationErrorMessage? registrationErrorMessage : ''}</span>
      </div>
      <div className="app__form-row">
        <button onClick={registerUser}>Registration</button>
      </div>
    </div>
  );
}

export default withRouter(Registration);