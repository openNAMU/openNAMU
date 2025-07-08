from .tool.func import *

async def view_xref(name = 'Test', xref_type = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check(name, 'render') == 1:
            return await re_error(conn, 0)

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        if xref_type == 1:
            div = '<a href="/xref_this/' + url_pas(name) + '">(' + await get_lang('link_in_this') + ')</a><hr class="main_hr">'
            data_sub = '(' + await get_lang('backlink') + ')'
        else:
            div = '<a href="/xref/' + url_pas(name) + '">(' + await get_lang('normal') + ')</a><hr class="main_hr">'
            data_sub = '(' + await get_lang('link_in_this') + ')'

        div += '<ul>'

        curs.execute(db_change('select data from other where name = "link_case_insensitive"'))
        db_data = curs.fetchall()
        link_case_insensitive = ' collate nocase' if db_data and db_data[0][0] != '' else ''

        if xref_type == 2:
            curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'link_count'"), [name])
            db_data = curs.fetchall()
            div += '<li>' + await get_lang('link_count') + ' : ' +  (db_data[0][0] if db_data else await get_lang('data_missing')) + '</li>'

        sql_insert = ['link', 'title'] if xref_type == 1 else ['title', 'link']
        curs.execute(db_change("select distinct " + sql_insert[0] + ", type from back where " + sql_insert[1] + " = ?" + link_case_insensitive + " and not type = 'no' and not type = 'nothing' order by type asc, " + sql_insert[0] + " asc limit ?, 50"), [name, sql_num])
        data_list = curs.fetchall()
        for data in data_list:
            div += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a>'

            if data[1]:
                div += ' (' + data[1] + ')'

            curs.execute(db_change("select title from back where title = ? and type = 'include'"), [data[0]])
            db_data = curs.fetchall()
            if db_data:
                div += ' <a class="opennamu_link_inter" href="/xref/' + url_pas(data[0]) + '">(' + await get_lang('backlink') + ')</a>'

            div += '</li>'

        div += '</ul>'
        
        if xref_type == 2:
            div += await get_next_page_bottom('/xref_this_page/{}/' + url_pas(name), num, data_list)
        else:
            div += await get_next_page_bottom('/xref_page/{}/' + url_pas(name), num, data_list)

        return easy_minify(flask.render_template(await skin_check(),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css([data_sub, 0])],
            data = div,
            menu = [['w/' + url_pas(name), await get_lang('return')], ['xref_reset/' + url_pas(name), await get_lang('reset_backlink')]]
        ))