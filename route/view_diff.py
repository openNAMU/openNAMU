from .tool.func import *

def view_diff(name = 'Test', num_a = 1, num_b = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        first = str(num_a)
        second = str(num_b)

        if acl_check(name, 'render') == 1:
            return re_error('/ban')

        curs.execute(db_change("select title from history where title = ? and (id = ? or id = ?) and hide = 'O'"), [name, first, second])
        if curs.fetchall() and admin_check(6) != 1:
            return re_error('/error/3')

        curs.execute(db_change("select data from history where id = ? and title = ?"), [first, name])
        first_raw_data = curs.fetchall()

        curs.execute(db_change("select data from history where id = ? and title = ?"), [second, name])
        second_raw_data = curs.fetchall()
        if first_raw_data and second_raw_data:
            first_raw_data = first_raw_data[0][0].replace('\r', '')
            second_raw_data = second_raw_data[0][0].replace('\r', '')

            if first_raw_data == second_raw_data:
                result = ''
            else:
                i = 1
                change_count = 0
                diff_data = diff_match_patch().diff_prettyHtml(
                    diff_match_patch().diff_main(first_raw_data, second_raw_data)
                )
                end_data = ''

                diff_data = diff_data.replace('&para;<br>', '\n')
                diff_data = diff_data.replace('<span>', '')
                diff_data = diff_data.replace('</span>', '')

                re_data = re.findall(r'(?:(?:(?:(?!\n).)*)(?:\n)|(?:(?:(?!\n).)+)$)', diff_data)
                for re_in_data in re_data:
                    change_find_start = len(re.findall(r'<(?:del|ins) ', re_in_data))
                    change_find_end = len(re.findall(r'<\/(?:del|ins)>', re_in_data))

                    change_count += (change_find_start - change_find_end)
                    if change_count != 0 or change_find_start != 0 or change_find_end != 0:
                        end_data += str(i) + ' : ' + re_in_data

                    i += 1

                result = '<pre>' + end_data + '</pre>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('compare') + ')', 0])],
                data = result,
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))
        else:
            return redirect('/history/' + url_pas(name))