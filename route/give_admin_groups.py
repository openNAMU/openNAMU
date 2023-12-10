from .tool.func import *

def give_admin_groups_2(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        acl_name_list = ['ban', 'nothing', 'toron', 'check', 'acl', 'hidel', 'give', 'owner']

        if flask.request.method == 'POST':
            if admin_check(None, 'auth list add (' + name + ')') != 1:
                return re_error('/error/3')
            elif name in get_default_admin_group():
                return re_error('/error/3')

            curs.execute(db_change("delete from alist where name = ?"), [name])
            for i in acl_name_list:
                if flask.request.form.get(i, 0) != 0:
                    curs.execute(db_change("insert into alist (name, acl) values (?, ?)"), [name, i])

            conn.commit()

            return redirect('/auth/list/add/' + url_pas(name))
        else:
            data = ''
            exist_list = ['', '', '', '', '', '', '', '']
            state = 'disabled' if admin_check() != 1 else ''
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

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('admin_group') + ')', 0])],
                data = '''
                    <form method="post">
                        ''' + data + '''
                        <hr class="main_hr">
                        <h2>''' + load_lang('explanation') + '''</h2>
                        <ul class="opennamu_ul">
                            <li style="margin-left: 20px;">owner : ''' + load_lang('owner_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">ban : ''' + load_lang('ban_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">toron : ''' + load_lang('discussion_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">check : ''' + load_lang('user_check_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">acl : ''' + load_lang('document_acl_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">hidel : ''' + load_lang('history_hide_authority') + '''</li>
                            <li style="margin-left: 40px; list-style: circle;">give : ''' + load_lang('authorization_authority') + '''</li>
                        </ul>
                        <hr class="main_hr">
                        <button ''' + state +  ''' type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', load_lang('return')]]
            ))