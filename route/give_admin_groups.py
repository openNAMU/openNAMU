from .tool.func import *

def give_admin_groups(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        acl_name_list = [
            [1, 'owner', get_lang(conn, 'owner_authority')],
            [2, 'admin', get_lang(conn, 'admin_authority')],
            [3, 'ban', get_lang(conn, 'ban_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'toron', get_lang(conn, 'discussion_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'check', get_lang(conn, 'user_check_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'acl', get_lang(conn, 'document_acl_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'hidel', get_lang(conn, 'history_hide_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'give', get_lang(conn, 'authorization_authority')],
            [4, '', get_lang(conn, 'admin_default_feature_authority'), True],
            [3, 'admin_default_feature', get_lang(conn, 'admin_default_feature_authority')],
            [4, 'user_name_bold', get_lang(conn, 'user_name_bold_authority')],
            [4, 'multiple_upload', get_lang(conn, 'multiple_upload_authority')],
            [4, 'slow_edit_pass', get_lang(conn, 'slow_edit_pass_authority')],
            [4, 'edit_bottom_compulsion_pass', get_lang(conn, 'edit_bottom_compulsion_pass_authority')],
            [4, '', get_lang(conn, 'user_authority'), True],
            [1, 'user', get_lang(conn, 'user_authority')],
            [2, 'captcha_pass', get_lang(conn, 'captcha_pass_authority')],
            [2, 'ip', get_lang(conn, 'ip_authority')],
            [3, 'view', get_lang(conn, 'view_authority')],
            [4, 'document', get_lang(conn, 'document_authority')],
            [5, 'edit', get_lang(conn, 'edit_authority')],
            [5, 'move', get_lang(conn, 'move_authority')],
            [5, 'new_make', get_lang(conn, 'new_make_authority')],
            [5, 'delete', get_lang(conn, 'delete_authority')],
            [4, 'edit_request', get_lang(conn, 'edit_request_authority')],
            [4, 'discuss', get_lang(conn, 'discuss_authority')],
            [4, 'upload', get_lang(conn, 'upload_authority')],
            [4, 'vote', get_lang(conn, 'vote_authority')],
            [4, 'captcha_one_check_five_pass', get_lang(conn, 'captcha_one_check_five_pass_authority')]
        ]

        if flask.request.method == 'POST':
            if admin_check(conn, None, 'auth list add (' + name + ')') != 1:
                return re_error(conn, '/error/3')
            elif name in get_default_admin_group():
                return re_error(conn, '/error/3')

            curs.execute(db_change("delete from alist where name = ?"), [name])
            for for_a in acl_name_list:
                if flask.request.form.get(for_a[1], 0) != 0:
                    curs.execute(db_change("insert into alist (name, acl) values (?, ?)"), [name, for_a[1]])

            return redirect(conn, '/auth/list/add/' + url_pas(name))
        else:
            state = 'disabled' if admin_check(conn) != 1 else ''
            state = 'disabled' if name in get_default_admin_group() else ''

            data = '<ul class="opennamu_ul">'
            for for_a in acl_name_list:
                curs.execute(db_change('select acl from alist where name = ?'), [name])
                acl_list = curs.fetchall()
                acl_list = [for_a[0] for for_a in acl_list]
                
                checked = ''
                if for_a[1] in acl_list:
                    checked = 'checked'

                choice = '<input type="checkbox" ' + state + ' name="' + for_a[1] + '" ' + checked + '> ' + for_a[2]
                if len(for_a) == 4:
                    choice = for_a[2]

                data += '' + \
                    '<li class="opennamu_list_1" style="margin-left: ' + str(int(for_a[0]) * 20) + 'px;">' + \
                        choice + \
                    '</li>'
                ''

            data += '</ul>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'admin_group') + ')', 0])],
                data = '''
                    <form method="post">
                        ''' + data + '''
                        <hr class="main_hr">
                        <button ''' + state +  ''' type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', get_lang(conn, 'return')]]
            ))