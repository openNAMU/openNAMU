from .tool.func import *

def list_old_page():
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        num = flask.request.args.get('num', '1')
        num = int(number_check(num))
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0
        
        curs.execute(db_change('select data from other where name = "count_all_title"'))
        if int(curs.fetchall()[0][0]) > 30000:
            return re_error('/error/25')
        
        div = '<ul class="inside_ul">'
        
        curs.execute(db_change('' + \
            'select h.title, max(h.date) from history as h where not (title like "user:%" or title like "category:%" or title like "file:%") and exists (select title from data where title = h.title) and not exists (select title from back where link = h.title and type = "redirect") group by h.title order by h.date asc limit ?, 50' + \
        ''), [sql_num])
        n_list = curs.fetchall()
        for data in n_list:
            div += '<li>' + data[1] + ' | <a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'
        
        div += '</ul>' + next_fix('/old_page?num=', num, n_list)
        
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('old_page'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))
