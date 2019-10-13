from .tool.func import *

def watch_list_2(conn):
    curs = conn.cursor()

    div = 'Limit : 10<hr class=\"main_hr\">'
    ip = ip_check()    

    if ip_or_user(ip) != 0:
        return redirect('/login')

    curs.execute("delete from scan where user = ? and title = ''", [ip])
    conn.commit()

    curs.execute("select title from scan where user = ?", [ip])
    data = curs.fetchall()
    for data_list in data:
        div += '<li><a href="/w/' + url_pas(data_list[0]) + '">' + data_list[0] + '</a> <a href="/watch_list/' + url_pas(data_list[0]) + '">(' + load_lang('delete') + ')</a></li>'

    if data:
        div = '<ul>' + div + '</ul><hr class=\"main_hr\">'

    div += '<a href="/manager/13">(' + load_lang('add') + ')</a>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('watchlist'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
    ))