from .tool.func import *

def user_ban_2(conn, name):
    curs = conn.cursor()

    if ip_or_user(name) == 0:
        curs.execute("select acl from user where id = ?", [name])
        user = curs.fetchall()
        if not user:
            return re_error('/error/2')

        if user and user[0][0] != 'user':
            if admin_check() != 1:
                return re_error('/error/4')

    if ban_check(ip = ip_check(), tool = 'login') == 1:
        return re_error('/ban')
                
    if flask.request.method == 'POST':
        if admin_check(1, 'ban (' + name + ')') != 1:
            return re_error('/error/3')

        if flask.request.form.get('limitless', '') == '':
            end = flask.request.form.get('second', '0')
        else:
            end = '0'

        ban_insert(name, end, flask.request.form.get('why', ''), flask.request.form.get('login', ''), ip_check())

        return redirect('/ban/' + url_pas(name))     
    else:
        if admin_check(1) != 1:
            return re_error('/error/3')

        curs.execute("select end, why from ban where block = ?", [name])
        end = curs.fetchall()
        if end:
            now = load_lang('release')

            if end[0][0] == '':
                data = '<ul><li>' + load_lang('limitless') + '</li>'
            else:
                data = '<ul><li>' + load_lang('period') + ' : ' + end[0][0] + '</li>'
                
            curs.execute("select block from ban where block = ? and login = 'O'", [name])
            if curs.fetchall():
                data += '<li>' + load_lang('login_able') + '</li>'

            if end[0][1] != '':
                data += '<li>' + load_lang('why') + ' : ' + end[0][1] + '</li></ul><hr class=\"main_hr\">'
            else:
                data += '</ul><hr class=\"main_hr\">'
        else:
            if re.search("^([0-9]{1,3}\.[0-9]{1,3})$", name):
                now = load_lang('band_ban')
            else:
                now = load_lang('ban')
                
            if ip_or_user(name) == 1:
                plus = '<input type="checkbox" name="login"> ' + load_lang('login_able') + '<hr class=\"main_hr\">'
            else:
                plus = ''

            data =  '''
                    <input placeholder="''' + load_lang('second') + '''" name="second" type="text">
                    <hr class=\"main_hr\">
                    <input type="checkbox" name="limitless"> ''' + load_lang('limitless') + '''
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="why" type="text">
                    <hr class=\"main_hr\">
                    ''' + plus

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + now + ')', 0])],
            data =  '''
                    <form method="post">
                        ''' + data + '''
                        <button type="submit">''' + now + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))   