from .set_mark.namumark import namumark, link_fix
from .set_mark.markdown import markdown

from .set_mark.tool import *

import re
import html
import sqlite3
import asyncio
import threading
import urllib.parse
import multiprocessing

def load_conn2(data):
    global conn
    global curs

    conn = data
    curs = conn.cursor()

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

def render_do(title, data, num, include):
    if num == 3:
        num = 1
        back_num = 3
    else:
        back_num = num

    curs.execute(db_change('select data from other where name = "markup"'))
    rep_data = curs.fetchall()
    if rep_data[0][0] == 'namumark':
        data = namumark(conn, data, title, num, include)
    elif rep_data[0][0] == 'raw':
        data = [data, '', []]
    else:
        data = ['', '', []]

    if num == 1:
        if data[2] == []:
            curs.execute(db_change("insert into back (title, link, type) values ('test', ?, 'nothing')"), [title])
        else:
            for data_in in data[2]:
                try:
                    curs.execute(db_change("insert into back (title, link, type) values (?, ?, ?)"), [data_in[1], data_in[0], data_in[2]])
                except:
                    pass

                curs.execute(db_change("delete from back where title = ? and type = 'no'"), [title])

        if back_num != 3:
            conn.commit()

    if num == 2:
        return [data[0], data[1]]
    else:
        return data[0] + data[1]