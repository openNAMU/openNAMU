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
                diff_data = diff_match_patch().diff_main(first_raw_data, second_raw_data)
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

                result = '<table style="width: 100%; white-space: pre-wrap;"><tr><td colspan="2">r' + first + ' ➤ r' + second + '</td></tr>'
                result += '<tr><td style="width: 40px; user-select: none;">'

                # 개행만 추가된 경우 조정 필요
                
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

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('compare') + ')', 0])],
                data = result,
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))
        else:
            return redirect('/history/' + url_pas(name))