from .tool.func import *
import pymysql

def main_error_404_2(conn):
    curs = conn.cursor()

    return redirect('/w/' + url_pas(wiki_set(2)))