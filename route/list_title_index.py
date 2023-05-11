from .tool.func import *

def list_title_index(num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        all_list = sql_num + 1
        data = ''

        curs.execute(db_change("select title from data order by title asc limit ?, 50"), [sql_num])
        title_list = curs.fetchall()
        if title_list:
            data += '<hr class="main_hr"><ul class="opennamu_ul">'

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
                    <ul class="opennamu_ul">
                        <li>''' + load_lang('all') + ' : ' + str(count_end[0]) + '''</li>
                    </ul>
                    <ul class="opennamu_ul">
                        <li>''' + load_lang('category') + ' : ' + str(count_end[1]) + '''</li>
                        <li>''' + load_lang('user_document') + ' : ' + str(count_end[2]) + '''</li>
                        <li>''' + load_lang('file') + ' : ' + str(count_end[3]) + '''</li>
                        <li>''' + load_lang('other') + ' : ' + str(count_end[4]) + '''</li>
                '''
            else:
                data += '''
                    </ul>
                    <ul class="opennamu_ul">
                        <li>''' + load_lang('all') + ' : ' + all_title[0][0] + '''</li>
                '''

        data += '</ul>' + next_fix('/list/document/all/', num, title_list)
        sub = ' (' + str(num) + ')'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('all_document_list'), wiki_set(), wiki_custom(), wiki_css([sub, 0])],
            data = data,
            menu = [['other', load_lang('return')]]
        ))