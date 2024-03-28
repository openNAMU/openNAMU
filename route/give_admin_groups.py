from .tool.func import *

def give_admin_groups_2(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        acl_name_list = ['ban', 'nothing', 'toron', 'check', 'acl', 'hidel', 'give', 'owner']

        if flask.request.method == 'POST':
            if admin_check(conn, None, 'auth list add (' + name + ')') != 1:
                return re_error(conn, '/error/3')
            elif name in get_default_admin_group():
                return re_error(conn, '/error/3')

            curs.execute(db_change("delete from alist where name = ?"), [name])
            for i in acl_name_list:
                if flask.request.form.get(i, 0) != 0:
                    curs.execute(db_change("insert into alist (name, acl) values (?, ?)"), [name, i])

            return redirect(conn, '/auth/list/add/' + url_pas(name))
        else:
            data = ''
            exist_list = ['', '', '', '', '', '', '', '']
            state = 'disabled' if admin_check(conn) != 1 else ''
            state = 'disabled' if name in get_default_admin_group() else ''

            curs.execute(db_change('select acl from alist where name = ?'), [name])
            acl_list = curs.fetchall()
            for go in acl_list:
                exist_list[acl_name_list.index(go[0])] = 'checked="checked"'

            for i in range(0, 8):
                if i != 1:
                    data += '' + \
                        '<input type="checkbox" ' + \
                                state + ' ' + \
                                'name="' + acl_name_list[i] + '" ' + \
                                exist_list[i] + '> ' + acl_name_list[i] + \
                        '<hr class="main_hr">' + \
                    ''

            data += ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'admin_group') + ')', 0])],
                data = '''
                    <form method="post">
                        ''' + data + '''
                        <hr class="main_hr">
                        <h2>''' + get_lang(conn, 'explanation') + '''</h2>
                        <ul class="opennamu_ul">
                            <li style="margin-left: 20px;">owner : ''' + get_lang(conn, 'owner_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">ban : ''' + get_lang(conn, 'ban_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">toron : ''' + get_lang(conn, 'discussion_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">check : ''' + get_lang(conn, 'user_check_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">acl : ''' + get_lang(conn, 'document_acl_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">hidel : ''' + get_lang(conn, 'history_hide_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">give : ''' + get_lang(conn, 'authorization_authority') + '''</li>
                        </ul>
                        <hr class="main_hr">
                        <button ''' + state +  ''' type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', get_lang(conn, 'return')]]
            ))