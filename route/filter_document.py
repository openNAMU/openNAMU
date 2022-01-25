from .tool.func import *

def filter_document():
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<table id="main_table_set">'
        div += '<tr id="main_table_top_tr">'
        div += '<td id="main_table_width">A</td>'
        div += '<td id="main_table_width">B</td>'
        div += '<td id="main_table_width">C</td>'
        div += '</tr>'

        admin = admin_check()
        title = load_lang('document_filter_list')
        
        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'document'"))
        db_data = curs.fetchall()
        for data in db_data:
            div += '<tr>'
            div += '<td>' + html.escape(data[0])
            if admin == 1:
                div += ' <a href="/filter/document/add/' + url_pas(data[0]) + '">(' + load_lang('edit') + ')</a>'
                div += ' <a href="/filter/document/del/' + url_pas(data[0]) + '">(' + load_lang('delete') + ')</a>'

            div += '</td>'
            div += '<td>' + html.escape(data[1]) + '</td>'
            div += '<td>' + html.escape(data[2]) + '</td>'
            div += '</tr>'

        div += '</table>'

        if admin == 1:
            div += '<hr class="main_hr">'
            div += '<a href="/filter/document/add">(' + load_lang('add') + ')</a>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [title, wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['manager/1', load_lang('return')]]
        ))