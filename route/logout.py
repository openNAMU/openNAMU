from .tool.func import *

def logout_2(conn):
    curs = conn.cursor()

    flask.session['state'] = 0
    flask.session.pop('id', None)

    return redirect('/user')