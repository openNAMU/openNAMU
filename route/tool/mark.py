from .namumark import namumark
from .tool import *

import re
import os
import html
import sqlite3
import asyncio
import threading
import urllib.parse

conn = ''
curs = ''

if os.path.exists('route/tool/custom.py'):
    from .custom import custom_mark
else:
    def custom_mark(conn, doc_data, doc_name, data_in):
        return [
            doc_data, 
            '', 
            []
        ]

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

def render_do(doc_name, doc_data, data_type, data_in):
    data_in = None if data_in == '' else data_in
    curs.execute(db_change('select data from other where name = "markup"'))
    rep_data = curs.fetchall()
    if rep_data[0][0] in ('namumark', 'js_onmark'):
        data_in = (data_in + '_') if data_in else ''
        data_end = [
            '<div class="render_content" id="' + data_in + 'render_content">' + html.escape(doc_data) + '</div>', 
            '''
                do_onmark_render(
                    test_mode = 0, 
                    name_id = "''' + data_in + '''render_content",
                    name_include = "''' + data_in + '''",
                    name_doc = "''' + doc_name.replace('"', '//"') + '''",
                );
            ''',
            []
        ]
    elif rep_data[0][0] == 'custom':
        data_end = custom_mark(
            conn, 
            doc_data, 
            doc_name, 
            data_in
        )
    else:
        data_end = [
            doc_data, 
            '', 
            []
        ]

    if data_type == 'backlink':
        if data_end[2] == []:
            curs.execute(db_change("insert into back (title, link, type) values ('test', ?, 'nothing')"), [doc_name])
        else:
            curs.executemany(db_change("insert into back (link, title, type) values (?, ?, ?)"), data_end[2])
            curs.execute(db_change("delete from back where title = ? and type = 'no'"), [doc_name])

        conn.commit()
    else:
        if data_type == 'api_view':
            return [
                data_end[0], 
                data_end[1]
            ]
        else:
            return data_end[0] + '<script>' + data_end[1] + '</script>'