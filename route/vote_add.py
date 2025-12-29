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
            acl_data = '<select class="__ON_SELECT__" name="acl_select">'
            acl_list = await get_acl_list()
            for data_list in acl_list:
                acl_data += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'

            acl_data += '</select>'

            return await render_template(
                await get_lang('add_vote'),
                '' + \
                    '<form method="post">' + \
                        '<input class="__ON_INPUT__" name="name" placeholder="' + await get_lang('name') + '">' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_100 __ON_TEXTAREA__" name="subject" placeholder="' + await get_lang('explanation') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_500 __ON_TEXTAREA__" name="data" placeholder="' + await get_lang('1_line_1_q') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<label><input class="__ON_INPUT__" type="checkbox" value="Y" name="open_select"> ' + await get_lang('open_vote') + '</label>' + \
                        '<h2>' + await get_lang('period') + '</h2>'
                        '<input class="__ON_INPUT__" type="date" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">' + \
                        '<hr class="main_hr">' + \
                        '<label><input class="__ON_INPUT__" type="checkbox" value="Y" name="limitless"> ' + await get_lang('limitless') + '</label>' + \
                        '<h2>' + await get_lang('acl') + '</h2>' + \
                        acl_data + ' <a href="/acl/TEST#exp">(' + await get_lang('explanation') + ')</a>' + \
                        '<hr class="main_hr">' + \
                        '<button class="__ON_BUTTON__" type="submit">' + await get_lang('send') + '</buttom>' + \
                    '</form>' + \
                '',
                0,
                [['vote', await get_lang('return')]]
            )
