from .tool.func import *

from .go_api_func_search import api_func_search

async def main_search_deep(name = 'Test', search_type = 'title', num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if name == '':
            return redirect(conn)

        if flask.request.method == 'POST':
            if search_type == 'title':
                return redirect(conn, '/search_page/1/' + url_pas(flask.request.form.get('search', 'test')))
            else:
                return redirect(conn, '/search_data_page/1/' + url_pas(flask.request.form.get('search', 'test')))
        else:
            div = '''
                <form method="post">
                    <input class="opennamu_width_200" name="search" value="''' + html.escape(name) + '''">
                    <button type="submit">''' + await get_lang('search') + '''</button>
                </form>
                <hr class="main_hr">
            '''

            if search_type == 'title':
                div += '<a href="/search_data_page/1/' + url_pas(name) + '">(' + await get_lang('search_document_data') + ')</a>'
            else:
                div += '<a href="/search_page/1/' + url_pas(name) + '">(' + await get_lang('search_document_name') + ')</a>'

            name_new = ''
            if re.search(r'^분류:', name):
                name_new = re.sub(r"^분류:", 'category:', name)
            elif re.search(r"^사용자:", name):
                name_new = re.sub(r"^사용자:", 'user:', name)
            elif re.search(r"^파일:", name):
                name_new = re.sub(r"^파일:", 'file:', name)

            if name_new != '':
                div += ' <a href="/search_page/1/' + url_pas(name_new) + '">(' + name_new + ')</a>'

            curs.execute(db_change("select title from data where title = ? collate nocase"), [name])
            link_id = '' if curs.fetchall() else 'class="opennamu_not_exist_link"'

            div += '''
                <ul>
                    <li>
                        ''' + await get_lang('go') + ''' : <a ''' + link_id + ' href="/w/' + url_pas(name) + '">' + html.escape(name) + '''</a>
                    </li>
                </ul>
                <ul>
            '''

            all_list = await api_func_search(name, search_type, num)
            for data in all_list:
                div += '<li><a href="/w/' + url_pas(data) + '">' + data + '</a></li>'

            div += '</ul>'
            
            if search_type == 'title':
                div += await get_next_page_bottom('/search_page/{}/' + url_pas(name), num, all_list)
            else:
                div += await get_next_page_bottom('/search_data_page/{}/' + url_pas(name), num, all_list)

            return easy_minify(flask.render_template(await skin_check(),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + await get_lang('search') + ')', 0])],
                data = div,
                menu = 0
            ))