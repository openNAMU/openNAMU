from .tool.func import *

def filter_all(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<table id="main_table_set">'
        div += '<tr id="main_table_top_tr">'

        div += '<td id="main_table_width">A</td>'
        div += '<td id="main_table_width">B</td>'
        div += '<td id="main_table_width">C</td>'

        div += '</tr>'

        admin = admin_check()

        if tool == 'inter_wiki':
            title = load_lang('interwiki_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'inter_wiki'"))
        elif tool == 'email_filter':
            title = load_lang('email_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'email'"))
        elif tool == 'name_filter':
            title = load_lang('id_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'name'"))
        elif tool == 'edit_filter':
            title = load_lang('edit_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'regex_filter'"))
        elif tool == 'file_filter':
            title = load_lang('file_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'file'"))
        elif tool == 'image_license':
            title = load_lang('image_license_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'image_license'"))
        elif tool == 'extension_filter':
            title = load_lang('extension_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'extension'"))
        elif tool == 'document':
            title = load_lang('document_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'document'"))
        elif tool == 'outer_link':
            title = load_lang('outer_link_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'outer_link'"))
        elif tool == 'template':
            title = load_lang('template_document_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'template'"))
        else:
            title = load_lang('edit_tool_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'edit_top'"))

        db_data = curs.fetchall()
        for data in db_data:
            div += '<tr>'
            div += '<td>'

            div += html.escape(data[0])
            if admin == 1:
                if tool in ('inter_wiki', 'outer_link', 'edit_filter', 'document', 'edit_top', 'template'):
                    div += ' <a href="/filter/' + tool + '/add/' + url_pas(data[0]) + '">(' + load_lang('edit') + ')</a>'
                    
                div += ' <a href="/filter/' + tool + '/del/' + url_pas(data[0]) + '">(' + load_lang('delete') + ')</a>'

            div += '</td>'

            if tool in ('inter_wiki', 'outer_link'):
                if tool == 'inter_wiki':
                    div += '<td><a class="opennamu_link_out" href="' + html.escape(data[1]) + '">' + html.escape(data[1]) + '</a></td>'
                else:
                    div += '<td>' + html.escape(data[1]) + '</td>'
                
                div += '<td>' + data[2] + '</td>'
            else:
                div += '<td>' + html.escape(data[1]) + '</td>'
                div += '<td>' + html.escape(data[2]) + '</td>'
            
            div += '</tr>'

        div += '</table>'

        if admin == 1:
            div += '<hr class="main_hr">'
            div += '<a href="/filter/' + tool + '/add">(' + load_lang('add') + ')</a>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [title, wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['manager/1', load_lang('return')]]
        ))