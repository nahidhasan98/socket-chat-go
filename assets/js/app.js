console.log("script linked successfully");

$(document).ready(() => {
    $('#message').keyup((event) => {
        if (event.shiftKey && event.keyCode == 13) {
            console.log("shift+enter pressed");
        } else if (event.keyCode == 13) {
            console.log($('#message').val().trim());
            $('#message').val("");
            //submit
        }
    });

    let socket = new WebSocket("ws://" + location.host + "/ws");
    console.log("Attempting Connection...");

    socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };

    socket.onclose = (event) => {
        console.log("Socket Closed Connection: ", event);
        socket.send("Client Closed!")
    };

    socket.onerror = (error) => {
        console.log("Socket Error: ", error);
    };
});

/* nginx reverse proxy setup for ws
location /ws {
    proxy_pass http://localhost:6001;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
    proxy_read_timeout 86400;
}
*/