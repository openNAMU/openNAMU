from .tool.func import *

def search_goto_2(name):
    

    if flask.request.form.get('search', None):
        data = flask.request.form.get('search', 'test')
    else:
        data = name

    sqlQuery("select title from data where title = ?", [data])
    t_data = sqlQuery("fetchall")
    if t_data:
        return redirect('/w/' + url_pas(data))
    else:
        return redirect('/search/' + url_pas(data))