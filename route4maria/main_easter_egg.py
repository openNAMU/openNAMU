from .tool.func import *
import pymysql

def main_easter_egg_2(conn):
    curs = conn.cursor()

    return easy_minify(flask.render_template(skin_check(), 
        imp = ['easter_egg.html', wiki_set(), custom(), other2([0, 0])],
        data = open('./views/easter_egg.html', 'r').read(),
        menu = 0
    ))