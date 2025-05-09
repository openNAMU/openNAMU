from .tool.func import *

async def filter_all(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = '<table id="main_table_set">'
        div += '<tr id="main_table_top_tr">'

        div += '<td id="main_table_width">A</td>'
        div += '<td id="main_table_width">B</td>'
        div += '<td id="main_table_width">C</td>'

        div += '</tr>'

        admin = await acl_check(tool = 'owner_auth')
        admin = 1 if admin == 0 else 0

        if tool == 'edit_filter':
            if await acl_check('', 'edit_filter_view', '', '') == 1:
                return await re_error(conn, 0)

        if tool == 'inter_wiki':
            title = get_lang(conn, 'interwiki_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'inter_wiki'"))
        elif tool == 'email_filter':
            title = get_lang(conn, 'email_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'email'"))
        elif tool == 'name_filter':
            title = get_lang(conn, 'id_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'name'"))
        elif tool == 'edit_filter':
            title = get_lang(conn, 'edit_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'regex_filter'"))
        elif tool == 'file_filter':
            title = get_lang(conn, 'file_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'file'"))
        elif tool == 'image_license':
            title = get_lang(conn, 'image_license_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'image_license'"))
        elif tool == 'extension_filter':
            title = get_lang(conn, 'extension_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'extension'"))
        elif tool == 'document':
            title = get_lang(conn, 'document_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'document'"))
        elif tool == 'outer_link':
            title = get_lang(conn, 'outer_link_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'outer_link'"))
        elif tool == 'template':
            title = get_lang(conn, 'template_document_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'template'"))
        else:
            title = get_lang(conn, 'edit_tool_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'edit_top'"))

        db_data = curs.fetchall()
        for data in db_data:
            div += '<tr>'
            div += '<td>'

            div += html.escape(data[0])
            if admin == 1:
                if tool in ('inter_wiki', 'outer_link', 'edit_filter', 'document', 'edit_top', 'template'):
                    div += ' <a href="/filter/' + tool + '/add/' + url_pas(data[0]) + '">(' + get_lang(conn, 'edit') + ')</a>'
                    
                div += ' <a href="/filter/' + tool + '/del/' + url_pas(data[0]) + '">(' + get_lang(conn, 'delete') + ')</a>'

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
            div += '<a href="/filter/' + tool + '/add">(' + get_lang(conn, 'add') + ')</a>'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [title, await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = div,
            menu = [['manager/1', get_lang(conn, 'return')]]
        ))