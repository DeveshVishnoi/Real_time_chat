const API_ENDPOINTS = "http://localhost:8080/api";

export async function loginHTTPRequest(email, password) {
    const response = await fetch(`${API_ENDPOINTS}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            email,
            password
        })
    });
    return await response.json();
}
export async function registerHTTPRequest(username,email, password) {
    const response = await fetch(`${API_ENDPOINTS}/reg`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username,
            email,
            password
        })
    });
    return await response.json();
}

export async function isEmailAvailableHTTPRequest(username) {
    const response = await fetch(`${API_ENDPOINTS}/isAvailable/${username}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    
    return await response.json();
}

export async function userSessionCheckHTTPRequest(username) {
    const response = await fetch(`${API_ENDPOINTS}/sessionStatus/${username}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });
    return await response.json();
}


export async function getConversationBetweenUsers(toUserID, fromUserID) {
    const response = await fetch(`${API_ENDPOINTS}/getConversation/${toUserID}/${fromUserID}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });
    return await response.json();
}
