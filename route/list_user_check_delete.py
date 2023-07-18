from .tool.func import *

def list_user_check_delete(name = None, ip = None, time = None, do_type = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/4')

        user_id = name
        user_ip = ip
        return_type = do_type

        if user_id and user_ip and time:
            if flask.request.method == 'POST':
                curs.execute(db_change("delete from ua_d where name = ? and ip = ? and today = ?"), [user_id, user_ip, time])
                conn.commit()

                return redirect('/list/user/check/' + url_pas(user_id if return_type == '0' else user_ip))
            else:
                return easy_minify(flask.render_template(skin_check(),
                    imp = [load_lang('check'), wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('delete') + ')', 0])],
                    data = '''
                        ''' + load_lang('name') + ''' : ''' + user_id + '''
                        <hr class="main_hr">
                        ''' + load_lang('ip') + ''' : ''' + user_ip + '''
                        <hr class="main_hr">
                        ''' + load_lang('time') + ''' : ''' + time + '''
                        <hr class="main_hr">
                        <form method="post">
                            <button type="submit">''' + load_lang('delete') + '''</button>
                        </form>
                    ''',
                    menu = [['check/' + url_pas(user_id if return_type == '0' else user_ip), load_lang('return')]]
                ))
        else:
            return redirect()