from .tool.func import *

def give_auth(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        owner_auth = admin_check()
        admin_auth = admin_check(7)

        curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [name])
        user_acl = curs.fetchall()
        if not user_acl:
            return re_error('/error/2')
        else:
            user_acl = user_acl[0][0]

        if owner_auth != 1:
            curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [user_acl])
            if curs.fetchall():
                return re_error('/error/3')

            if ip_check() == name:
                return re_error('/error/3')

        if flask.request.method == 'POST':
            if admin_check(7, 'admin (' + name + ')') != 1:
                return re_error('/error/3')

            select_data = flask.request.form.get('select', 'X')
            if select_data == 'X':
                select_data = 'user'

            curs.execute(db_change('select name from alist where name = ? and acl = "owner"'), [select_data])
            if owner_auth != 1 and curs.fetchall():
                return re_error('/error/3')

            curs.execute(db_change("update user_set set data = ? where id = ? and name = 'acl'"), [select_data, name])
            curs.execute(db_change('delete from user_set where name = "auth_date" and id = ?'), [name])

            time_limit = flask.request.form.get('date', '')
            if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                curs.execute(db_change("insert into user_set (id, name, data) values (?, 'auth_date', ?)"), [name, time_limit])

            conn.commit()

            return redirect('/auth/give/' + url_pas(name))
        else:
            if admin_auth != 1:
                return re_error('/error/3')

            div = '<option value="X">' + load_lang('normal') + '</option>'
            div += '<option value="ban">' + load_lang('ban') + '</option>'

            curs.execute(db_change('select distinct name from alist order by name asc'))
            for data in curs.fetchall():
                if user_acl == data[0]:
                    div = '<option value="' + data[0] + '">' + data[0] + '</option>' + div
                else:
                    div += '<option value="' + data[0] + '">' + data[0] + '</option>'
                    
            date_value = ''
            
            curs.execute(db_change('select data from user_set where name = "auth_date" and id = ?'), [name])
            db_data = curs.fetchall()
            if db_data:
                date_value = db_data[0][0]

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('authorize') + ')', 0])],
                data =  '''
                    <form method="post">
                        <div id="opennamu_get_user_info">''' + html.escape(name) + '''</div>
                        <hr class="main_hr">
                        <select name="select">''' + div + '''</select>
                        <hr class="main_hr">
                        <input type="date" value="''' + date_value + '''" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['manager', load_lang('return')]]
            ))