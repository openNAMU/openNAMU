from .tool.func import *

def give_user_check_delete_2(conn):
    curs = conn.cursor()
    
    if admin_check() != 1:
        return re_error('/error/4')

    user_id = flask.request.args.get('name', None)
    user_ip = flask.request.args.get('ip', None)
    
    time = flask.request.args.get('time', None)
    time_set = re.search(r'([0-9]{4})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})', time)
    if not time_set:
        return redirect()
    
    time_set = time_set.groups()
    time = time_set[0] + '-' + time_set[1] + '-' + time_set[2] + ' '
    time += time_set[3] + ':' + time_set[4] + ':' + time_set[5]
    
    return_type = flask.request.args.get('return_type', '1')
        
    if user_id and user_ip and time:
        if flask.request.method == 'POST':
            curs.execute(db_change("delete from ua_d where name = ? and ip = ? and today = ?"), [user_id, user_ip, time])
            conn.commit()
            
            return redirect('/check/' + url_pas(user_id if return_type == '0' else user_ip))
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