from .tool.func import *

async def list_admin_auth_use(arg_num = 1, arg_search = 'normal'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        if flask.request.method == 'POST':
            return redirect(conn, '/list/admin/auth_use_page/1/' + url_pas(flask.request.form.get('search', 'normal')))
        else:
            arg_search = 'normal' if arg_search == '' else arg_search
            
            if arg_search == 'normal':
                curs.execute(db_change("select who, what, time from re_admin order by time desc limit ?, 50"), [sql_num])
            else:
                curs.execute(db_change("select who, what, time from re_admin where what like ? order by time desc limit ?, 50"), [arg_search + "%", sql_num])

            list_data = '<ul>'

            get_list = curs.fetchall()
            for data in get_list:
                do_data = data[1]

                if ip_or_user(data[0]) != 0:
                    curs.execute(db_change("select data from other where name = 'ip_view'"))
                    db_data = curs.fetchall()
                    ip_view = db_data[0][0] if db_data else ''
                    ip_view = '' if await acl_check(tool = 'ban_auth') != 1 else ip_view
                    
                    if ip_view != '':
                        do_data = do_data.split(' ')
                        do_data = do_data[0] if do_data[0] in ['ban'] else data[1]

                list_data += '<li>' + await ip_pas(data[0]) + ' | ' + html.escape(do_data) + ' | ' + data[2] + '</li>'

            list_data += '</ul>'
            list_data += get_next_page_bottom(conn, '/list/admin/auth_use_page/{}/' + url_pas(arg_search), arg_num, get_list)

            arg_search = html.escape(arg_search) if arg_search != 'normal' else ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'authority_use_list'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input class="opennamu_width_200" name="search" placeholder="''' + get_lang(conn, 'start_with_search') + '''" value="''' + arg_search + '''">
                        <button type="submit">''' + get_lang(conn, 'search') + '''</button>
                    </form>
                    <hr class="main_hr">
                ''' + list_data,
                menu = [['other', get_lang(conn, 'return')]]
            ))