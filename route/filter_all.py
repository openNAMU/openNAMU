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
            title = await get_lang('interwiki_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'inter_wiki'"))
        elif tool == 'email_filter':
            title = await get_lang('email_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'email'"))
        elif tool == 'name_filter':
            title = await get_lang('id_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'name'"))
        elif tool == 'edit_filter':
            title = await get_lang('edit_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'regex_filter'"))
        elif tool == 'file_filter':
            title = await get_lang('file_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'file'"))
        elif tool == 'image_license':
            title = await get_lang('image_license_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'image_license'"))
        elif tool == 'extension_filter':
            title = await get_lang('extension_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'extension'"))
        elif tool == 'document':
            title = await get_lang('document_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'document'"))
        elif tool == 'outer_link':
            title = await get_lang('outer_link_filter_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'outer_link'"))
        elif tool == 'template':
            title = await get_lang('template_document_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'template'"))
        else:
            title = await get_lang('edit_tool_list')
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'edit_top'"))

        db_data = curs.fetchall()
        for data in db_data:
            div += '<tr>'
            div += '<td>'

            div += html.escape(data[0])
            if admin == 1:
                if tool in ('inter_wiki', 'outer_link', 'edit_filter', 'document', 'edit_top', 'template'):
                    div += ' <a href="/filter/' + tool + '/add/' + url_pas(data[0]) + '">(' + await get_lang('edit') + ')</a>'
                    
                div += ' <a href="/filter/' + tool + '/del/' + url_pas(data[0]) + '">(' + await get_lang('delete') + ')</a>'

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
            div += '<a href="/filter/' + tool + '/add">(' + await get_lang('add') + ')</a>'

        return easy_minify(flask.render_template(await skin_check(),
            imp = [title, await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['manager/1', await get_lang('return')]]
        ))