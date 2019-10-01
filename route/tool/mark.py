from .set_mark.namu import namu
from .set_mark.markdown import markdown

import re
import html
import json
import pymysql
import sqlite3
# import psycopg2 as pg2
import urllib.parse
import threading
import multiprocessing

def load_conn2(data):
    global conn
    global curs

    conn = data
    curs = conn.cursor()

def q_mariadb2(query, *arg):
    if arg:
        qarg = tuple(arg[0])
        sql = query.replace("%", "%%")
        sql = sql.replace("?", "%s")
        sql = sql.replace("'", '"')

        return curs.execute(sql, qarg)
    elif not arg:
        sql = query.replace("%", "%%")
        sql = sql.replace("'", '"')
        sql = sql.replace("random()", "RAND()")

        return curs.execute(sql)

def q_sqlite2(query, *arg):
    if arg:
        qarg = arg[0]
        return curs.execute(query, qarg)
    elif not arg:
        return curs.execute(query)

def sqlQuery2(query, *arg):
    db_type = {}
    try:
        open('DB_Type.json', mode='r', encoding='utf-8')
    except FileNotFoundError:
        print("error!")

    with open("DB_Type.json") as fileRead:
        db_type = json.load(fileRead)

    if db_type["DBMS"] == "mariadb":
        if query == "fetchall":
            return curs.fetchall()
        elif query == "commit":
            pass
        else:
            return q_mariadb2(query, *arg)
    elif db_type["DBMS"] == "sqlite":
        if query == "fetchall":
            return curs.fetchall()
        elif query == "commit":
            return conn.commit()
        else:
            return q_sqlite2(query, *arg)
    else:
        print("DBMS Type Error")

def send_parser(data):
    if not re.search('^<br>$', data):
        data = html.escape(data)
        
        javascript = re.compile('javascript:', re.I)
        
        data = javascript.sub('', data)

        while 1:
            re_data = re.search('&lt;a(?: (?:(?:(?!&gt;).)*))?&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', data)
            if re_data:
                re_data = re_data.groups()[0]

                data = re.sub('&lt;a(?: (?:(?:(?!&gt;).)*))?&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="/w/' + urllib.parse.quote(re_data).replace('/','%2F') + '">' + re_data + '</a>', data, 1)
            else:
                break
        
    return data
    
def plusing(data):
    for data_in in data:
        sqlQuery2("select title from back where title = ? and link = ? and type = ?", [data_in[1], data_in[0], data_in[2]])
        if not sqlQuery2("fetchall"):
            sqlQuery2("insert into back (title, link, type) values (?, ?, ?)", [data_in[1], data_in[0], data_in[2]])

def namumark(conn, title = '', data = None, num = 0):
    sqlQuery2('select data from other where name = "markup"')
    rep_data = sqlQuery2("fetchall")
    if rep_data[0][0] == 'namumark':
        data = namu(conn, data, title, num)
    elif rep_data[0][0] == 'markdown':
        data = markdown(conn, data, title, num)
    else:
        data = ['', '', []]

    if num == 1:
        data_num = len(data[2]) 
        data_in_num = int(data_num / multiprocessing.cpu_count())
        data_in = []

        for i in range(multiprocessing.cpu_count()):
            if i != multiprocessing.cpu_count() - 1:
                data_in += [data[2][data_in_num * i:data_in_num * (i + 1)]]
            else:
                data_in += [data[2][data_in_num * i:]]

        for data_in_for in data_in:
            thread_start = threading.Thread(target = plusing, args = [data_in_for])
            thread_start.start()
            thread_start.join()
        
        sqlQuery2("commit")
        
    if num == 2:
        return [data[0], data[1]]
    else:
        return data[0] + data[1]