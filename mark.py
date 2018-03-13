from set_mark.start import start

import re
import html
import sqlite3
from urllib import parse
import time
import threading

def send_parser(data):
    data = html.escape(data)
    
    javascript = re.compile('javascript:', re.I)
    
    data = javascript.sub('', data)
    data = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + parse.quote('\g<in>').replace('/','%2F') + '">\g<in></a>', data)  
    
    return data
    
def plusing(conn, name, link, backtype):
    curs = conn.cursor()
    curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
    if not curs.fetchall():
        curs.execute("insert into back (title, link, type) values (?, ?, ?)", [link, name, backtype])

def namumark(conn, title, data, num):
    data = start(conn, data, title)
    if num == 1:        
        for back_data in data[2]:
            thread_start = threading.Thread(target = plusing, args = [conn, back_data[0], back_data[1], back_data[2]])
            thread_start.start()
            thread_start.join()

        conn.commit()
        
    return data[0] + data[1]