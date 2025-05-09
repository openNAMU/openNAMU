from .tool.func import *

async def list_user_check_delete(name = None, ip = None, time = None, do_type = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 4)

        user_id = name
        user_ip = ip
        return_type = do_type

        if user_id and user_ip and time:
            if flask.request.method == 'POST':
                curs.execute(db_change("delete from ua_d where name = ? and ip = ? and today = ?"), [user_id, user_ip, time])

                return redirect(conn, '/list/user/check/' + url_pas(user_id if return_type == '0' else user_ip))
            else:
                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'check'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'delete') + ')', 0])],
                    data = '''
                        ''' + get_lang(conn, 'name') + ''' : ''' + user_id + '''
                        <hr class="main_hr">
                        ''' + get_lang(conn, 'ip') + ''' : ''' + user_ip + '''
                        <hr class="main_hr">
                        ''' + get_lang(conn, 'time') + ''' : ''' + time + '''
                        <hr class="main_hr">
                        <form method="post">
                            <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                        </form>
                    ''',
                    menu = [['check/' + url_pas(user_id if return_type == '0' else user_ip), get_lang(conn, 'return')]]
                ))
        else:
            return redirect(conn)