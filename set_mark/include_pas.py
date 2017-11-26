from . import html_pas
from . import link
from . import mid_pas
from . import toc_pas
import sqlite3
from urllib import parse
import re

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))
    
def include_pas(conn, data, title, in_c, num, toc_y, fol_num):
    curs = conn.cursor()

    category = ''
    backlink = []
    
    include = re.compile("\[include\(((?:(?!\)\]|,).)*)((?:(?:,\s?(?:(?!\)\]).)*))+)?\)\]")
    m = include.findall(data)
    for results in m:
        if(results[0] == title):
            data = include.sub("<b>" + results[0] + "</b>", data, 1)
        else:
            curs.execute("select data from data where title = ?", [results[0]])
            in_con = curs.fetchall()
            
            backlink += [[title, results[0], 'include']]
            if(in_con):                        
                in_data = in_con[0][0]
                in_data = include.sub("", in_data)
                in_data = re.sub("\n", "\r\n", re.sub("\r\n", "\n", in_data))
                in_data = html_pas.html_pas(in_data)
                
                var_d = mid_pas.mid_pas(in_data, fol_num, 1, in_c)
                var_d2 = link.link(conn, title, var_d[0], 0, category, backlink)

                in_data = var_d2[0]
                category = var_d2[1]
                fol_num = var_d[1]
                
                if(results[1]):
                    a = results[1]
                    while(1):
                        g = re.search("([^= ,]*)\=([^,]*)", a)
                        if(g):
                            result = g.groups()
                            in_data = re.sub("@" + result[0] + "@", result[1], in_data)
                            a = re.sub("([^= ,]*)\=([^,]*)", "", a, 1)
                        else:
                            break       

                in_data = toc_pas.toc_pas(in_data, results[0], num, toc_y)
                            
                data = include.sub('\n<nobr><a id="include_link" href="/w/' + url_pas(results[0]) + '">[' + results[0] + ' 이동]</a><div>' + in_data + '</div><nobr>\n', data, 1)
            else:
                data = include.sub("<a class=\"not_thing\" href=\"/w/" + url_pas(results[0]) + "\">" + results[0] + "</a>", data, 1)

    return([data, category, fol_num, backlink])