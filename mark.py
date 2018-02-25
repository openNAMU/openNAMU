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

    data = re.sub('&lt;a href="(?:[^"]*)"&gt;(?P<in>(?:(?!&lt;).)*)&lt;\/a&gt;', '<a href="' + parse.quote('\g<in>').replace('/','%2F') + '">\g<in></a>', d)  
    return data
    
def plusing(conn, name, link, backtype):
    curs = conn.cursor()

    curs.execute("select title from back where title = ? and link = ? and type = ?", [link, name, backtype])
    if not curs.fetchall():
        curs.execute("insert into back (title, link, type) values (?, ?, ?)", [link, name, backtype])

def namumark(conn, title, data, num, in_c, toc_y):
    data = start(conn, data, title)

    data += '<script>function folding(num, test = 0) { var fol = document.getElementById(\'folding_\' + num); if(fol.style.display == \'inline-block\' || fol.style.display == \'block\') { fol.style.display = \'none\'; } else { if(num % 3 == 0 && test != 1) { fol.style.display = \'block\'; } else { fol.style.display = \'inline-block\'; } } } </script>'
    
    # if num == 1:        
    #     for d4 in backlink:
    #         t = threading.Thread(target = plusing, args = [conn, d4[0], d4[1], d4[2]])
    #         t.start()
    #         t.join()
    #
    #     conn.commit()
        
    return data