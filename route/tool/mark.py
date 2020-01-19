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

def plusing(data_in):
    try:
        curs.execute(db_change("insert into back (title, link, type) values (?, ?, ?)"), [data_in[1], data_in[0], data_in[2]])
    except:
        pass

def render_do(title, data, num, include):
    curs.execute(db_change('select data from other where name = "markup"'))
    rep_data = curs.fetchall()
    if rep_data[0][0] == 'namumark':
        data = namumark(conn, data, title, num, include)
    elif rep_data[0][0] == 'js_namumark':
        data = [
            '<div id="render_contect">' + html.escape(data) + '</div>',
            '<script>render_namumark("render_contect")</script>',
            []
        ]
    elif rep_data[0][0] == 'markdown':
        data = markdown(conn, data, title, num)
    elif rep_data[0][0] == 'raw':
        data = [data, '', []]
    else:
        data = ['', '', []]

    if num == 1:
        for data_in in data[2]:
            plusing(data_in)

        conn.commit()

    if num == 2:
        return [data[0], data[1]]
    else:
        return data[0] + data[1]