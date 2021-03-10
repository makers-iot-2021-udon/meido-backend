import socketio
from flask import Flask
sio = socketio.Server(async_mode='threading')
app = Flask(__name__)

@sio.event
def connect(sid, environ, auth):
    print('connect ', sid)

@sio.event
def disconnect(sid):
    print('disconnect ', sid)

app.wsgi_app = socketio.WSGIApp(sio, app.wsgi_app)

# ... Socket.IO and Flask handler functions ...

if __name__ == '__main__':
    app.run(threaded=True)
# ... Socket.IO and Flask handler functions ...

# # wrap with a WSGI application
# app = socketio.WSGIApp(sio)
