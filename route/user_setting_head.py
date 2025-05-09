from .tool.func import *

async def user_setting_head(skin_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()

        skin_name_org = skin_name
        if skin_name != '':
            skin_name = '_' + skin_name
    
        if flask.request.method == 'POST':
            get_data = flask.request.form.get('content', '')
            if ip_or_user(ip) == 0:
                curs.execute(db_change("select id from user_set where id = ? and name = ?"), [ip, 'custom_css' + skin_name])
                if curs.fetchall():
                    curs.execute(db_change("update user_set set data = ? where id = ? and name = ?"), [get_data, ip, 'custom_css' + skin_name])
                else:
                    curs.execute(db_change("insert into user_set (id, name, data) values (?, ?, ?)"), [ip, 'custom_css' + skin_name, get_data])
        
            flask.session['head' + skin_name] = get_data

            if skin_name_org != '':
                return redirect(conn, '/change/head/' + skin_name_org)
            else:
                return redirect(conn, '/change/head')
        else:
            if ip_or_user(ip) == 0:
                start = ''

                curs.execute(db_change("select data from user_set where id = ? and name = ?"), [ip, 'custom_css' + skin_name])
                head_data = curs.fetchall()
                data = head_data[0][0] if head_data else ''
            else:
                start = '' + \
                    '<span>' + get_lang(conn, 'user_head_warning') + '</span>' + \
                    '<hr class="main_hr">' + \
                ''
                data = flask.session['head' + skin_name] if 'head' + skin_name in flask.session else ''

            start += '' + \
                '<span>' + \
                    '&lt;style&gt;CSS&lt;/style&gt;' + \
                    '<br>' + \
                    '&lt;script&gt;JS&lt;/script&gt;' + \
                '</span>' + \
                '<hr class="main_hr">' + \
            ''

            if skin_name == '':
                sub_name = ''
            else:
                sub_name = ' (' + skin_name_org + ')'

            start = '' + \
                '<a href="/change/head">(' + get_lang(conn, 'all') + ')</a> ' + \
                ' '.join(['<a href="/change/head/' + url_pas(i) + '">(' + html.escape(i) + ')</a>' for i in load_skin(conn, '', 1)]) + \
                '<hr class="main_hr">' + \
                start + \
            ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, data = 'user_head', safe = 1), await wiki_set(), await wiki_custom(conn), wiki_css(['(HTML)' + sub_name, 0])],
                data = start + '''
                    <form method="post">
                        <textarea class="opennamu_textarea_500" cols="100" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        ''' + get_lang(conn, 'user_css_warning') + ''' : <a href="/change/head_reset">/change/head_reset</a>
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['change', get_lang(conn, 'return')]]
            ))