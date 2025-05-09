from .tool.func import *

def view_diff_do(first_raw_data, second_raw_data, first, second):
    if first_raw_data == second_raw_data:
        result = ''
    else:
        dmp = diff_match_patch()

        diff_data = dmp.diff_main(first_raw_data, second_raw_data)
        dmp.diff_cleanupSemantic(diff_data)

        diff_data += [[0, '\n']]
        
        diff_data_2 = []
        temp_list = []
        line = 1
        line_change = 0
        for for_a in diff_data:
            line_split = re.findall(r'(.*\n)|(.+$)', for_a[1])
            if line_split:
                for for_b in line_split:
                    if for_b[0] != '':
                        if for_a[0] != 0:
                            line_change = 1
                        
                        temp_list += [[line, for_a[0], for_b[0].replace('\n', '')]]

                        if line_change == 1:
                            diff_data_2 += temp_list
                        
                        temp_list = []
                        line_change = 0
                        line += 1
                    else:
                        if for_a[0] != 0:
                            line_change = 1

                        temp_list += [[line, for_a[0], for_b[1]]]
            else:
                if for_a[0] != 0:
                    line_change = 1

                temp_list += [[line, for_a[0], for_a[1]]]

        result = '<table style="width: 100%; white-space: pre-wrap;"><tr><td colspan="2">' + first + ' ➤ ' + second + '</td></tr>'
        result += '<tr><td style="width: 40px; user-select: none;">'

        line = 0
        for for_a in diff_data_2:
            if line == 0:
                line = for_a[0]
                result += str(line) + '</td><td>'
            else:
                if line != for_a[0]:
                    line = for_a[0]
                    result += '</td></tr><tr><td style="width: 40px; user-select: none;">' + str(line) + '</td><td>'

            if for_a[1] == 1:
                result += '<span class="opennamu_diff_green">' + html.escape(for_a[2]) + '</span>'
            elif for_a[1] == 0:
                result += html.escape(for_a[2])
            else:
                result += '<span class="opennamu_diff_red">' + html.escape(for_a[2]) + '</span>'

        result += '</td></tr></table>'

    return result

async def view_diff(name = 'Test', num_a = 1, num_b = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        first = str(num_a)
        second = str(num_b)

        if await acl_check(name, 'render') == 1:
            return await re_error(conn, 0)

        curs.execute(db_change("select title from history where title = ? and (id = ? or id = ?) and hide = 'O'"), [name, first, second])
        if curs.fetchall() and await acl_check(tool = 'hidel_auth') == 1:
            return await re_error(conn, 3)

        curs.execute(db_change("select data from history where id = ? and title = ?"), [first, name])
        first_raw_data = curs.fetchall()

        curs.execute(db_change("select data from history where id = ? and title = ?"), [second, name])
        second_raw_data = curs.fetchall()
        
        if first_raw_data and second_raw_data:
            first_raw_data = first_raw_data[0][0].replace('\r', '')
            second_raw_data = second_raw_data[0][0].replace('\r', '')

            result = view_diff_do(first_raw_data, second_raw_data, 'r' + first, 'r' + second)

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'compare') + ')', 0])],
                data = result,
                menu = [['history/' + url_pas(name), get_lang(conn, 'return')]]
            ))
        else:
            return redirect(conn, '/history/' + url_pas(name))