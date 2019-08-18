from .tool.func import *
import pymysql

def user_custom_head_view_2(conn):
    curs = conn.cursor()

    ip = ip_check()

    if flask.request.method == 'POST':
        if custom()[2] != 0:
            curs.execute("select user from custom where user = %s", [ip + ' (head)'])
            if curs.fetchall():
                curs.execute("update custom set css = %s where user = %s", [flask.request.form.get('content', None), ip + ' (head)'])
            else:
                curs.execute("insert into custom (user, css) values (%s, %s)", [ip + ' (head)', flask.request.form.get('content', None)])
            
            conn.commit()

        flask.session['head'] = flask.request.form.get('content', None)

        return redirect('/user')
    else:
        if custom()[2] != 0:
            start = ''

            curs.execute("select css from custom where user = %s", [ip + ' (head)'])
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