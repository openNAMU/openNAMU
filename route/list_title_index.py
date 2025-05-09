from .tool.func import *

async def list_title_index(num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        all_list = sql_num + 1
        data = ''

        curs.execute(db_change("select title from data order by title asc limit ?, 50"), [sql_num])
        title_list = curs.fetchall()
        if title_list:
            data += '<hr class="main_hr"><ul>'

        for list_data in title_list:
            data += '<li>' + str(all_list) + '. <a href="/w/' + url_pas(list_data[0]) + '">' + html.escape(list_data[0]) + '</a></li>'
            all_list += 1

        if num == 1:
            count_end = []

            curs.execute(db_change('select data from other where name = "count_all_title"'))
            all_title = curs.fetchall()
            if int(all_title[0][0]) < 30000:
                count_end += [int(all_title[0][0])]

                sql_list = ['category:', 'user:', 'file:']
                for sql in sql_list:
                    curs.execute(db_change("select count(*) from data where title like ?"), [sql + '%'])
                    count = curs.fetchall()
                    if count:
                        count_end += [int(count[0][0])]
                    else:
                        count_end += [0]

                count_end += [count_end[0] - count_end[1]  - count_end[2]  - count_end[3]]

                data += '''
                    </ul>
                    <ul>
                        <li>''' + get_lang(conn, 'all') + ' : ' + str(count_end[0]) + '''</li>
                    </ul>
                    <ul>
                        <li>''' + get_lang(conn, 'category') + ' : ' + str(count_end[1]) + '''</li>
                        <li>''' + get_lang(conn, 'user_document') + ' : ' + str(count_end[2]) + '''</li>
                        <li>''' + get_lang(conn, 'file') + ' : ' + str(count_end[3]) + '''</li>
                        <li>''' + get_lang(conn, 'other') + ' : ' + str(count_end[4]) + '''</li>
                '''
            else:
                data += '''
                    </ul>
                    <ul>
                        <li>''' + get_lang(conn, 'all') + ' : ' + all_title[0][0] + '''</li>
                '''

        data += '</ul>' + get_next_page_bottom(conn, '/list/document/all/{}', num, title_list)
        sub = ' (' + str(num) + ')'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'all_document_list'), await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
            data = data,
            menu = [['other', get_lang(conn, 'return')]]
        ))