from .set_mark.namumark import namumark
from .set_mark.markdown import markdown

from .set_mark.tool import *

import re
import os
import html
import sqlite3
import asyncio
import threading
import urllib.parse

conn = ''
curs = ''

if os.path.exists('route/tool/set_mark/custom.py'):
    from .set_mark.custom import custom_mark
else:
    def custom_mark(conn, data, title, num, include):
        return [data, '', []]

def load_conn2(data):
    global conn
    global curs

    conn = data
    curs = conn.cursor()

def send_parser(data):
    if not re.search(r'^<br>$', data):
        data = html.escape(data)

        data = re.sub(r'javascript:', '', data, flags = re.I)
        data = data.replace('&lt;br&gt;', '')

    link_re = re.compile(r'&lt;a(?: (?:(?:(?!&gt;).)*))?&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;')
    link_data = link_re.findall(data)
    for i in link_data:
        data = link_re.sub('<a href="/w/' + urllib.parse.quote(i).replace('/','%2F') + '">' + i + '</a>', data, 1)

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
        data = namumark(conn, data, title, include)
    elif rep_data[0][0] == 'markdown':
        data = markdown(conn, data, title, include)
    elif rep_data[0][0] == 'custom':
        data = custom_mark(conn, data, title, include)
    elif rep_data[0][0] == 'raw':
        data = [data, '', []]
    else:
        data = ['', '', []]

    if num == 1:
        if data[2] == []:
            curs.execute(db_change("insert into back (title, link, type) values ('test', ?, 'nothing')"), [title])
        else:
            curs.executemany(db_change("insert into back (link, title, type) values (?, ?, ?)"), data[2])
            curs.execute(db_change("delete from back where title = ? and type = 'no'"), [title])

        if back_num != 3:
            conn.commit()

    if num == 2:
        return [data[0], data[1]]
    else:
        return data[0] + '<script>' + data[1] + '</script>'