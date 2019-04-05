from .tool.func import *

def logout_2(conn):
    curs = conn.cursor()

    flask.session.pop('state', None)
    flask.session.pop('id', None)

    return redirect('/user')