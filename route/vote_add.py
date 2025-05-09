from .tool.func import *

async def vote_add():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'vote') == 1:
            return await re_error(conn, 0)

        if flask.request.method == 'POST':
            vote_data = flask.request.form.get('data', 'test\ntest_2')
            if vote_data.count('\n') < 1:
                return await re_error(conn, 0)

            curs.execute(db_change('select id from vote where not type = "option" order by id + 0 desc limit 1'))
            id_data = curs.fetchall()
            id_data = str((int(id_data[0][0]) + 1) if id_data else 1)

            if flask.request.form.get('open_select', 'N') == 'Y':
                open_data = 'open'
            else:
                open_data = 'n_open'

            curs.execute(db_change("insert into vote (name, id, subject, data, user, type, acl) values (?, ?, ?, ?, '', ?, ?)"), [
                flask.request.form.get('name', 'test'),
                id_data,
                flask.request.form.get('subject', 'test'),
                flask.request.form.get('data', 'test'),
                open_data,
                flask.request.form.get('acl_select', '')
            ])
            curs.execute(db_change("insert into vote (name, id, subject, data, user, type, acl) values ('open_user', ?, '', ?, '', 'option', '')"), [
                id_data,
                ip_check()
            ])
            
            time_limitless = flask.request.form.get('limitless', '')
            if time_limitless == '':
                time_limit = flask.request.form.get('date', '')
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                    curs.execute(db_change("insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')"), [
                        id_data,
                        time_limit
                    ])

            return redirect(conn, '/vote')
        else:
            acl_data = '<select name="acl_select">'
            acl_list = await get_acl_list()
            for data_list in acl_list:
                acl_data += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'

            acl_data += '</select>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'add_vote'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '' + \
                    '<form method="post">' + \
                        '<input name="name" placeholder="' + get_lang(conn, 'name') + '">' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_100" name="subject" placeholder="' + get_lang(conn, 'explanation') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_500" name="data" placeholder="' + get_lang(conn, '1_line_1_q') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<label><input type="checkbox" value="Y" name="open_select"> ' + get_lang(conn, 'open_vote') + '</label>' + \
                        '<h2>' + get_lang(conn, 'period') + '</h2>'
                        '<input type="date" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">' + \
                        '<hr class="main_hr">' + \
                        '<label><input type="checkbox" value="Y" name="limitless"> ' + get_lang(conn, 'limitless') + '</label>' + \
                        '<h2>' + get_lang(conn, 'acl') + '</h2>' + \
                        acl_data + ' <a href="/acl/TEST#exp">(' + get_lang(conn, 'explanation') + ')</a>' + \
                        '<hr class="main_hr">' + \
                        '<button type="submit">' + get_lang(conn, 'send') + '</buttom>' + \
                    '</form>' + \
                '',
                menu = [['vote', get_lang(conn, 'return')]]
            ))