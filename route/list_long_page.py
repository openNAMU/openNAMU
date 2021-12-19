from .tool.func import *

def list_long_page_2(conn, tool):
    curs = conn.cursor()

    curs.execute(db_change('select data from other where name = "count_all_title"'))
    if int(curs.fetchall()[0][0]) > 30000:
        return re_error('/error/25')

    div = '<ul class="inside_ul">'
    select_data = 'desc' if tool == 'long_page' else 'asc'
    title = 'long_page' if tool == 'long_page' else 'short_page'

    curs.execute(db_change("select title, length(data) from data order by length(data) " + select_data + " limit 50"))
    for data in curs.fetchall():
        div += '<li>' + str(data[1]) + ' : <a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

    div += '</ul>'

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang(title), wiki_set(), wiki_custom(), wiki_css([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))
