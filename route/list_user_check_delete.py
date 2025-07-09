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
                return easy_minify(flask.render_template(await skin_check(),
                    imp = [await get_lang('check'), await wiki_set(), await wiki_custom(), wiki_css(['(' + await get_lang('delete') + ')', 0])],
                    data = '''
                        ''' + await get_lang('name') + ''' : ''' + user_id + '''
                        <hr class="main_hr">
                        ''' + await get_lang('ip') + ''' : ''' + user_ip + '''
                        <hr class="main_hr">
                        ''' + await get_lang('time') + ''' : ''' + time + '''
                        <hr class="main_hr">
                        <form method="post">
                            <button type="submit">''' + await get_lang('delete') + '''</button>
                        </form>
                    ''',
                    menu = [['check/' + url_pas(user_id if return_type == '0' else user_ip), await get_lang('return')]]
                ))
        else:
            return redirect(conn)