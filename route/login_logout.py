from .tool.func import *

def login_logout_2():
    

    flask.session.pop('state', None)
    flask.session.pop('id', None)

    return redirect('/user')