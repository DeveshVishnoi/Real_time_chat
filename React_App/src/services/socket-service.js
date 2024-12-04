const events = require('events');

const CHAT_SERVER_ENDPOINT = "https://real-time-chat-wmhv.onrender.com/api";
let webSocketConnection = null;

export const eventEmitter = new events.EventEmitter();

export function connectToWebSocket(userID) {

    
    if (userID === "" && userID === null && userID === undefined) {
        return {
            message: "You need User ID to connect to the Chat server",
            webSocketConnection: null
        }
    } else if (!window["WebSocket"]) {
        return {
            message: "Your Browser doesn't support Web Sockets",
            webSocketConnection: null
        }
    }
    if (window["WebSocket"]) {
        webSocketConnection = new WebSocket("ws://" + CHAT_SERVER_ENDPOINT + "/ws/" + userID);
        return {
            message: "You are connected to Chat Server",
            webSocketConnection
        }
    }
}

export function sendWebSocketMessage(messagePayload) {
    if (webSocketConnection === null) {
      return;
    }
    webSocketConnection.send(
      JSON.stringify({
        message_type: 'message',
        message: messagePayload
      })
    );
}

export function emitLogoutEvent(userID) {
    if (webSocketConnection === null) {
        return;
    }
    webSocketConnection.close();
}

export function listenToWebSocketEvents() {

    if (webSocketConnection === null) {
        return;
    }

    webSocketConnection.onclose = (event) => {
        eventEmitter.emit('disconnect', event);
    };

    webSocketConnection.onmessage = (event) => {


        console.log("Event data : ", event.data);
        

        try {
            const socketPayload = JSON.parse(event.data);

            console.log(socketPayload);
            
            switch (socketPayload.message_type) {
                case 'chatlist-response':
                    if (!socketPayload.message) {
                        return
                    }
                    eventEmitter.emit(
                      'chatlist-response',
                      socketPayload.message
                    );

                    break;

                case 'disconnect':
                    if (!socketPayload.message) {
                        return
                    }
                    eventEmitter.emit(
                      'chatlist-response',
                      socketPayload.message
                    );

                    break;

                case 'message-response':

                    if (!socketPayload.message) {
                        return
                    }

                    eventEmitter.emit('message-response', socketPayload.message);
                    break;

                default:
                    break;
            }
        } catch (error) {
            console.log(error)
            console.warn('Something went wrong while decoding the Message Payload')
        }
    };
}
