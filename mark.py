from set_mark.namu import namu

import re
import html
import sqlite3
import urllib.parse
import time
import threading
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
        data = re.sub('&lt;a href=&quot;(?:(?:(?!&quot;).)*)&quot;&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + urllib.parse.quote('\g<in>').replace('/','%2F') + '">\g<in></a>', data)
    
    return data
    
def plusing(data):
    for data_in in data:
        curs.execute("select title from back where title = ? and link = ? and type = ?", [data_in[1], data_in[0], data_in[2]])
        if not curs.fetchall():
            curs.execute("insert into back (title, link, type) values (?, ?, ?)", [data_in[1], data_in[0], data_in[2]])

def namumark(title = '', data = '', num = 0):
    if not data == '':
        data = namu(conn, data, title, num)

        if num == 1:
            data_num = len(data[2]) 
            data_in_num = int(data_num / 8)
            data_in = []

            for i in range(8):
                if not i == 7:
                    data_in += [data[2][data_in_num * i:data_in_num * (i + 1)]]
                else:
                    data_in += [data[2][data_in_num * i:]]

            for data_in_for in data_in:
                thread_start = threading.Thread(target = plusing, args = [data_in_for])
                thread_start.start()
                thread_start.join()
            
            conn.commit()
            
        return data[0] + data[1]
    else:
        return '404'