from .tool.func import *

def user_custom_head_view_2(conn):
    curs = conn.cursor()

    ip = ip_check()

    if flask.request.method == 'POST':
        if ip_or_user(ip) == 0:
            curs.execute(db_change("select user from custom where user = ?"), [ip + ' (head)'])
            if curs.fetchall():
                curs.execute(db_change("update custom set css = ? where user = ?"), [flask.request.form.get('content', None), ip + ' (head)'])
            else:
                curs.execute(db_change("insert into custom (user, css) values (?, ?)"), [ip + ' (head)', flask.request.form.get('content', None)])

            conn.commit()

        flask.session['head'] = flask.request.form.get('content', None)

        return redirect('/custom_head')
    else:
        if ip_or_user(ip) == 0:
            start = ''

            curs.execute(db_change("select css from custom where user = ?"), [ip + ' (head)'])
            head_data = curs.fetchall()
            if head_data:
                data = head_data[0][0]
            else:
                data = ''
        else:
            start = '<span>' + load_lang('user_head_warring') + '</span><hr class=\"main_hr\">'

            if 'head' in flask.session:
                data = flask.session['head']
            else:
                data = ''

        start += '<span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span><hr class=\"main_hr\">'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang(data = 'user_head', safe = 1), wiki_set(), custom(), other2([0, 0])],
            data =  start + '''
                    <form method="post">
                        <textarea rows="25" cols="100" name="content">''' + data + '''</textarea>
                        <hr class=\"main_hr\">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))