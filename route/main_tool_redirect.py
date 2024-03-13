from .tool.func import *

def main_tool_redirect(num = 1, add_2 = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        title_list = {
            0 : [get_lang(conn, 'document_name'), '/acl', get_lang(conn, 'acl')],
            1 : [0, '/list/user/check', get_lang(conn, 'check')],
            2 : [get_lang(conn, 'file_name'), '/filter/file_filter/add', get_lang(conn, 'file_filter_add')],
            3 : [0, '/auth/give', get_lang(conn, 'authorize')],
            4 : [0, '/user', get_lang(conn, 'user_tool')],
            6 : [get_lang(conn, 'name'), '/auth/list/add', get_lang(conn, 'add_admin_group')],
            7 : [get_lang(conn, 'name'), '/filter/edit_filter/add', get_lang(conn, 'edit_filter_add')],
            8 : [get_lang(conn, 'document_name'), '/search', get_lang(conn, 'search')],
            9 : [0, '/block_log/user', get_lang(conn, 'blocked_user')],
            10 : [0, '/block_log/admin', get_lang(conn, 'blocked_admin')],
            11 : [get_lang(conn, 'document_name'), '/watch_list', get_lang(conn, 'add_watchlist')],
            12 : [get_lang(conn, 'compare_target'), '/list/user/check', get_lang(conn, 'compare_target')],
            13 : [get_lang(conn, 'document_name'), '/edit', get_lang(conn, 'load')],
            14 : [get_lang(conn, 'document_name'), '/star_doc', get_lang(conn, 'add_star_doc')],
            15 : [get_lang(conn, 'name_or_ip_or_regex'), '/auth/give/ban', get_lang(conn, 'release')],
            16 : [0, '/auth/give/fix', get_lang(conn, 'user_fix')],
        }
        
        if num == 1:
            return redirect(conn, '/manager')
        elif num - 1 <= len(title_list):
            num -= 2

            add_1 = flask.request.form.get('name', 'test')
            if flask.request.method == 'POST':
                if add_2 != '':
                    if num != 12:
                        flask.session['edit_load_document'] = add_1
                        return redirect(conn, '/edit_from/' + url_pas(add_2))
                    else:
                        return redirect(conn, title_list[num][1] + '/' + url_pas(add_2) + '/normal/1/' + url_pas(add_1))
                elif flask.request.form.get('regex', '') != '':
                    return redirect(conn, '/auth/give/ban_regex/' + url_pas(add_1))
                else:
                    return redirect(conn, title_list[num][1] + '/' + url_pas(add_1))
            else:
                if title_list[num][0] == 0:
                    placeholder = get_lang(conn, 'user_name')
                else:
                    placeholder = title_list[num][0]

                plus = ''
                if num == 15:
                    plus = '<input type="checkbox" name="regex"> ' + get_lang(conn, 'regex') + '<hr class="main_hr">'

                top_plus = ''
                if num == 13:
                    curs.execute(db_change("select html, plus from html_filter where kind = 'template'"))
                    db_data = curs.fetchall()
                    for for_a in db_data:
                        top_plus += '' + \
                            '<a href="javascript:opennamu_insert_v(\'data_field\', \'' + get_tool_js_safe(for_a[0]) + '\')">' + html.escape(for_a[0]) + '</a> : ' + html.escape(for_a[1]) + \
                            '<hr class="main_hr">' + \
                        ''

                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [title_list[num][2], wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            ''' + top_plus + '''
                            <input placeholder="''' + placeholder + '''" id="data_field" name="name" type="text">
                            <hr class="main_hr">
                            ''' + plus + '''
                            <button type="submit">''' + get_lang(conn, 'go') + '''</button>
                        </form>
                    ''',
                    menu = [['manager', get_lang(conn, 'return')]]
                ))
        else:
            return redirect(conn)