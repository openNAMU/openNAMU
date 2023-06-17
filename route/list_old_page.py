from .tool.func import *

def list_old_page(num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        # 리다이렉트 구분도 넣을 예정
        # 그 전에 로직 개편하고
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0
        
        curs.execute(db_change('select data from other where name = "count_all_title"'))
        if int(curs.fetchall()[0][0]) > 30000:
            return re_error('/error/25')
        
        div = '<ul class="opennamu_ul">'
        
        curs.execute(db_change("select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' order by set_data asc limit ?, 50"), [sql_num])
        n_list = curs.fetchall()
        for data in n_list:
            div += '<li>' + data[1] + ' | <a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'
        
        div += '</ul>' + next_fix('/list/document/old/', num, n_list)
        
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('old_page'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))
