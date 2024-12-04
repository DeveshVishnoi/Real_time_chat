import React, {useState} from 'react';
import { withRouter } from 'react-router-dom';

import { loginHTTPRequest } from "./../../../services/api-service";
import { setItemInLS } from "./../../../services/storage-service";

import './login.css'

function Login(props) {

  const [loginErrorMessage, setErrorMessage] = useState(null);
  const [email, updateEmail] = useState(null);
  const [password, updatePassword] = useState(null);


  const handleEmailChanges = (event) => {
    updateEmail(event.target.value)
  }

  const handlePasswordChange = (event) => {
    updatePassword(event.target.value)
  }

  const loginUser = async () => {
    props.displayPageLoader(true);
    const userDetails = await loginHTTPRequest(email, password);
    props.displayPageLoader(false);

    
    if (userDetails.statusCode === 200) {
      setItemInLS('userDetails', userDetails.data)
      props.history.push(`/home`)
    } else {
      setErrorMessage(userDetails.message);
    }
  };

  return (
    <div className="app__login-container">
      <div className="app__form-row">
        <label>Email:</label>
        <input type="email" className="email" onChange={handleEmailChanges} />
      </div>
      <div className="app__form-row">
        <label>Password:</label>
        <input type="password" className="password" onChange={handlePasswordChange} />
      </div>
      <div className="app__form-row">
        <span className="error-message">{loginErrorMessage? loginErrorMessage : ''}</span>
      </div>
      <div className="app__form-row">
        <button onClick={loginUser}>Login</button>
      </div>
    </div>
  );
}

export default withRouter(Login);