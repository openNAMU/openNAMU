from .tool.func import *

def main_search_deep(name = 'Test', search_type = 'title', num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if name == '':
            return redirect()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        if flask.request.method == 'POST':
            if search_type == 'title':
                return redirect('/search/1/' + url_pas(flask.request.form.get('search', 'test')))
            else:
                return redirect('/search_data/1/' + url_pas(flask.request.form.get('search', 'test')))
        else:
            div = '''
                <form method="post">
                    <input class="opennamu_width_200" name="search" value="''' + html.escape(name) + '''">
                    <button type="submit">''' + load_lang('search') + '''</button>
                </form>
                <hr class="main_hr">
            '''

            if search_type == 'title':
                div += '<a href="/search_data/1/' + url_pas(name) + '">(' + load_lang('search_document_data') + ')</a>'
            else:
                div += '<a href="/search/1/' + url_pas(name) + '">(' + load_lang('search_document_name') + ')</a>'

            name_new = ''
            if re.search(r'^분류:', name):
                name_new = re.sub(r"^분류:", 'category:', name)
            elif re.search(r"^사용자:", name):
                name_new = re.sub(r"^사용자:", 'user:', name)
            elif re.search(r"^파일:", name):
                name_new = re.sub(r"^파일:", 'file:', name)

            if name_new != '':
                div += ' <a href="/search/1/' + url_pas(name_new) + '">(' + name_new + ')</a>'

            curs.execute(db_change("select title from data where title = ? collate nocase"), [name])
            link_id = '' if curs.fetchall() else 'class="opennamu_not_exist_link"'

            div += '''
                <ul class="opennamu_ul">
                    <li>
                        <a ''' + link_id + ' href="/w/' + url_pas(name) + '">' + html.escape(name) + '''</a>
                    </li>
                </ul>
                <hr class="main_hr">
                <ul class="opennamu_ul">
            '''

            if search_type == 'title':
                curs.execute(db_change("select title from data where title like ? collate nocase order by title limit ?, 50"),
                    ['%' + name + '%', sql_num]
                )
            else:
                curs.execute(db_change("select title from data where data like ? collate nocase order by title limit ?, 50"),
                    ['%' + name + '%', sql_num]
                )

            all_list = curs.fetchall()
            for data in all_list:
                div += '<li><a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a></li>'

            div += '</ul>'
            
            if search_type == 'title':
                div += get_next_page_bottom('/search/{}/' + url_pas(name), num, all_list)
            else:
                div += get_next_page_bottom('/search_data/{}/' + url_pas(name), num, all_list)

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('search') + ')', 0])],
                data = div,
                menu = 0
            ))