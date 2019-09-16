from .tool.func import *

def func_title_random_2():
    

    sqlQuery("select title from data order by random() limit 1")
    data = sqlQuery("fetchall")
    if data:
        return redirect('/w/' + url_pas(data[0][0]))
    else:
        return redirect()