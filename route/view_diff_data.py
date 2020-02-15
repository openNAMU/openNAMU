from .tool.func import *

def view_diff_data_2(conn, name):
    curs = conn.cursor()

    if acl_check(name, 'render') == 1:
        return re_error('/ban')

    first = number_check(flask.request.args.get('first', '1'))
    second = number_check(flask.request.args.get('second', '1'))

    curs.execute(db_change("select title from history where title = ? and (id = ? or id = ?) and hide = 'O'"), [name, first, second])
    if curs.fetchall() and admin_check(6) != 1:
        return re_error('/error/3')

    curs.execute(db_change("select data from history where id = ? and title = ?"), [first, name])
    first_raw_data = curs.fetchall()
    if first_raw_data:
        curs.execute(db_change("select data from history where id = ? and title = ?"), [second, name])
        second_raw_data = curs.fetchall()
        if second_raw_data:
            if first == second:
                result = ''
            elif first_raw_data == second_raw_data:
                result = ''
            else:
                i = 2
                include_ins = 0
                diff_data = '1 : ' + diff_match_patch().diff_prettyHtml(diff_match_patch().diff_main(first_raw_data[0][0], second_raw_data[0][0]))

                re_data = re.findall('((?:(?!&para;<br>).)*)(?:&para;<br>|$)', diff_data)
                for re_in_data in re_data:
                    if re.search('(<ins |<del )', re_in_data) and re.search('(<\/ins>|<\/del>)', re_in_data):
                        re_data[i - 2] = str(i) + ' : ' + re_data[i - 2] + '\n'
                        include_ins = 0
                    elif re.search('(<\/ins>|<\/del>)', re_in_data):
                        re_data[i - 2] = str(i) + ' : ' + re_data[i - 2] + '\n'
                        include_ins = 0
                    elif re.search('(<ins |<del )', re_in_data) or include_ins == 1:
                        re_data[i - 2] = str(i) + ' : ' + re_data[i - 2] + '\n'
                        include_ins = 1
                    else:
                        re_data[i - 2] = ''
                        include_ins = 0

                    i += 1
                    
                result = '<pre>' + ''.join(re_data) + '</pre>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), custom(), other2([' (' + load_lang('compare') + ')', 0])],
                data = result,
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))

    return redirect('/history/' + url_pas(name))