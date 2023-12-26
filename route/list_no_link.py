from .tool.func import *

def list_no_link(num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0
        
        div = '<ul class="opennamu_ul">'
        
        curs.execute(db_change("select doc_name, set_data from data_set where set_name = 'link_count' and doc_rev = '' and set_data = '0' limit ?, 50"), [sql_num])
        n_list = curs.fetchall()
        for data in n_list:
            div += '<li>'
            div += data[1] + ' | <a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a>'
            
            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'doc_type'"), [data[0]])
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                div += ' | ' + db_data[0][0]

            div += '</li>'
        
        div += '</ul>' + next_fix('/list/document/no_link/', num, n_list)
        
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('no_link_document_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))
