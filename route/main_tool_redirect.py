from .tool.func import *

def main_tool_redirect(num = 1, add_2 = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        title_list = {
            0 : [load_lang('document_name'), '/acl', load_lang('acl')],
            1 : [0, '/list/user/check', load_lang('check')],
            2 : [load_lang('file_name'), '/filter/file_filter/add', load_lang('file_filter_add')],
            3 : [0, '/auth/give', load_lang('authorize')],
            4 : [0, '/user', load_lang('user_tool')],
            6 : [load_lang('name'), '/auth/list/add', load_lang('add_admin_group')],
            7 : [load_lang('name'), '/filter/edit_filter/add', load_lang('edit_filter_add')],
            8 : [load_lang('document_name'), '/search', load_lang('search')],
            9 : [0, '/block_log/user', load_lang('blocked_user')],
            10 : [0, '/block_log/admin', load_lang('blocked_admin')],
            11 : [load_lang('document_name'), '/watch_list', load_lang('add_watchlist')],
            12 : [load_lang('compare_target'), '/list/user/check', load_lang('compare_target')],
            13 : [load_lang('document_name'), '/edit', load_lang('load')],
            14 : [load_lang('document_name'), '/star_doc', load_lang('add_star_doc')],
            15 : [load_lang('name_or_ip_or_regex'), '/auth/give/ban', load_lang('release')],
            16 : [0, '/auth/give/fix', load_lang('user_fix')],
        }
        
        if num == 1:
            return redirect('/manager')
        elif num - 1 <= len(title_list):
            num -= 2

            add_1 = flask.request.form.get('name', 'test')
            if flask.request.method == 'POST':
                if add_2 != '':
                    if num != 12:
                        flask.session['edit_load_document'] = add_1
                        return redirect('/edit_from/' + url_pas(add_2))
                    else:
                        return redirect(title_list[num][1] + '/' + url_pas(add_2) + '/normal/1/' + url_pas(add_1))
                elif flask.request.form.get('regex', '') != '':
                    return redirect('/auth/give/ban_regex/' + url_pas(add_1))
                else:
                    return redirect(title_list[num][1] + '/' + url_pas(add_1))
            else:
                if title_list[num][0] == 0:
                    placeholder = load_lang('user_name')
                else:
                    placeholder = title_list[num][0]

                plus = ''
                if num == 15:
                    plus = '<input type="checkbox" name="regex"> ' + load_lang('regex') + '<hr class="main_hr">'

                top_plus = ''
                if num == 13:
                    curs.execute(db_change("select html, plus from html_filter where kind = 'template'"))
                    db_data = curs.fetchall()
                    for for_a in db_data:
                        top_plus += '' + \
                            '<a href="javascript:opennamu_insert_v(\'data_field\', \'' + get_tool_js_safe(for_a[0]) + '\')">' + html.escape(for_a[0]) + '</a> : ' + html.escape(for_a[1]) + \
                            '<hr class="main_hr">' + \
                        ''

                return easy_minify(flask.render_template(skin_check(),
                    imp = [title_list[num][2], wiki_set(), wiki_custom(), wiki_css([0, 0])],
                    data = '''
                        <form method="post">
                            ''' + top_plus + '''
                            <input placeholder="''' + placeholder + '''" id="data_field" name="name" type="text">
                            <hr class="main_hr">
                            ''' + plus + '''
                            <button type="submit">''' + load_lang('go') + '''</button>
                        </form>
                    ''',
                    menu = [['manager', load_lang('return')]]
                ))
        else:
            return redirect()