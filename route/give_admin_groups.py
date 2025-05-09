from .tool.func import *

async def give_admin_groups(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        acl_name_list = [
            [1, 'owner', get_lang(conn, 'owner_authority')],
                [2, '', get_lang(conn, 'all_function_authority')],
                [2, 'admin', get_lang(conn, 'admin_authority')],
                    [3, 'ban', get_lang(conn, 'ban_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'toron', get_lang(conn, 'discussion_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'check', get_lang(conn, 'user_analyze_authority')],
                        [4, 'view_user_watchlist', get_lang(conn, 'view_user_watchlist_authority')],
                        [4, '', get_lang(conn, 'user_check_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'acl', get_lang(conn, 'document_acl_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'hidel', get_lang(conn, 'history_hide_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'give', get_lang(conn, 'authorization_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'bbs', get_lang(conn, 'bbs_management_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'vote_fix', get_lang(conn, 'vote_management_authority')],
                        [4, '', get_lang(conn, 'admin_default_feature_authority')],
                    [3, 'admin_default_feature', get_lang(conn, 'admin_default_feature_authority')],
                        [4, 'doc_watch_list_view', get_lang(conn, 'doc_watch_list_view_authority')],
                        [4, 'treat_as_admin', get_lang(conn, 'treat_as_admin_authority')],
                        [4, 'view_hide_user_name', get_lang(conn, 'view_hide_user_name_authority')],
                        [4, 'user_name_bold', get_lang(conn, 'user_name_bold_authority')],
                        [4, 'multiple_upload', get_lang(conn, 'multiple_upload_authority')],
                        [4, 'slow_edit_pass', get_lang(conn, 'slow_edit_pass_authority')],
                        [4, 'edit_bottom_compulsion_pass', get_lang(conn, 'edit_bottom_compulsion_pass_authority')],
                        [4, 'edit_filter_pass', get_lang(conn, 'edit_filter_pass_authority')],
                        [4, '', get_lang(conn, 'user_authority')],
            [1, 'user', get_lang(conn, 'user_authority')],
                [2, 'captcha_pass', get_lang(conn, 'captcha_pass_authority')],
                [2, 'ip', get_lang(conn, 'ip_authority')],
                    [3, 'document', get_lang(conn, 'document_authority')],
                        [4, 'edit', get_lang(conn, 'edit_authority')],
                            [5, '', get_lang(conn, 'view_authority')],
                        [4, 'edit_request', get_lang(conn, 'edit_request_authority')],
                            [5, '', get_lang(conn, 'view_authority')],
                        [4, 'move', get_lang(conn, 'move_authority')],
                            [5, '', get_lang(conn, 'view_authority')],
                        [4, 'new_make', get_lang(conn, 'new_make_authority')],
                            [5, '', get_lang(conn, 'view_authority')],
                        [4, 'delete', get_lang(conn, 'delete_authority')],
                            [5, '', get_lang(conn, 'view_authority')],
                        [4, 'view', get_lang(conn, 'view_authority')],
                    [3, 'discuss', get_lang(conn, 'discuss_authority')],
                        [4, 'discuss_make_new_thread', get_lang(conn, 'discuss_make_new_thread_authority')],
                            [5, '', get_lang(conn, 'discuss_view_authority')],
                        [4, 'discuss_view', get_lang(conn, 'discuss_view_authority')],
                    [3, 'upload', get_lang(conn, 'upload_authority')],
                    [3, 'vote', get_lang(conn, 'vote_authority')],
                    [3, 'bbs_use', get_lang(conn, 'bbs_authority')],
                        [4, 'bbs_edit', get_lang(conn, 'bbs_edit_authority')],
                            [5, '', get_lang(conn, 'bbs_view_authority')],
                        [4, 'bbs_comment', get_lang(conn, 'bbs_comment_authority')],
                            [5, '', get_lang(conn, 'bbs_view_authority')],
                        [4, 'bbs_view', get_lang(conn, 'bbs_view_authority')],
                    [3, 'captcha_one_check_five_pass', get_lang(conn, 'captcha_one_check_five_pass_authority')],
                    [3, 'edit_filter_view', get_lang(conn, 'edit_filter_view_authority')],
                    [3, 'nothing', get_lang(conn, 'nothing_authority')]
        ]

        if html.escape(name) != name:
            return await re_error(conn, 48)

        if flask.request.method == 'POST':
            if await acl_check(tool = 'owner_auth', memo = 'auth list add (' + name + ')') == 1:
                return await re_error(conn, 3)

            curs.execute(db_change("delete from alist where name = ?"), [name])
            for for_a in acl_name_list:
                if flask.request.form.get(for_a[1], 0) != 0:
                    curs.execute(db_change("insert into alist (name, acl) values (?, ?)"), [name, for_a[1]])

            curs.execute(db_change("insert into alist (name, acl) values (?, 'nothing')"), [name])

            return redirect(conn, '/auth/list/add/' + url_pas(name))
        else:
            state = 'disabled' if await acl_check('', 'owner_auth', '', '') == 1 else ''

            curs.execute(db_change('select acl from alist where name = ?'), [name])
            acl_list = curs.fetchall()
            acl_list = [for_b[0] for for_b in acl_list]

            data = '<ul>'
            for for_a in acl_name_list:                
                checked = ''
                if for_a[1] in acl_list:
                    checked = 'checked'
                    
                choice = '<label><input type="checkbox" ' + state + ' name="' + for_a[1] + '" ' + checked + '> ' + for_a[2] + '</label>'
                if for_a[1] == '':
                    choice = for_a[2]

                data += '' + \
                    '<li class="opennamu_list_1" style="margin-left: ' + str((int(for_a[0]) - 1) * 20) + 'px;">' + \
                        choice + \
                    '</li>'
                ''

            data += '</ul>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'admin_group') + ')', 0])],
                data = '''
                    <form method="post">
                        ''' + data + '''
                        <hr class="main_hr">
                        <button ''' + state +  ''' type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', get_lang(conn, 'return')]]
            ))