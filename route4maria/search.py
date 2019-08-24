from .tool.func import *
import pymysql

def search_2(conn):
    curs = conn.cursor()

    return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))
