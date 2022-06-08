from .tool.func import *

def user_setting_head():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
    
        if flask.request.method == 'POST':
            get_data = flask.request.form.get('content', '')
            if ip_or_user(ip) == 0:
                curs.execute(db_change("select id from user_set where id = ? and name = 'custom_css'"), [ip])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = 'custom_css'"), [get_data, ip])
                else:
                    curs.execute(db_change("insert into user_set (id, name, data) values (?, 'custom_css', ?)"), [ip, get_data])

                conn.commit()
        
            flask.session['head'] = get_data

            return redirect('/change/head')
        else:
            if ip_or_user(ip) == 0:
                start = ''

                curs.execute(db_change("select data from user_set where id = ? and name = 'custom_css'"), [ip])
                head_data = curs.fetchall()
                data = head_data[0][0] if head_data else ''
            else:
                start = '' + \
                    '<span>' + load_lang('user_head_warning') + '</span>' + \
                    '<hr class="main_hr">' + \
                ''
                data = flask.session['head'] if 'head' in flask.session else ''

            start += '' + \
                '<span>' + \
                    '&lt;style&gt;CSS&lt;/style&gt;' + \
                    '<br>' + \
                    '&lt;script&gt;JS&lt;/script&gt;' + \
                '</span>' + \
                '<hr class="main_hr">' + \
            ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang(data = 'user_head', safe = 1), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = start + '''
                    <form method="post">
                        <textarea rows="25" cols="100" name="content">''' + data + '''</textarea>
                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))