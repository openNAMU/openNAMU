from .tool.func import *

def vote_add():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
            vote_data = flask.request.form.get('data', 'test\ntest_2')
            if vote_data.count('\n') < 1:
                return re_error('/ban')

            curs.execute(db_change('select id from vote where not type = "option" order by id + 0 desc limit 1'))
            id_data = curs.fetchall()
            id_data = str((int(id_data[0][0]) + 1) if id_data else 1)

            admin_check(None, 'add vote ' + id_data)

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
            
            time_limitless = flask.request.form.get('limitless', '')
            if time_limitless == '':
                time_limit = flask.request.form.get('date', '')
                if re.search(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$', time_limit):
                    curs.execute(db_change("insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')"), [
                        id_data,
                        time_limit
                    ])
                    
            conn.commit()

            return redirect('/vote')
        else:
            acl_data = '<select name="acl_select">'
            acl_list = get_acl_list()
            for data_list in acl_list:
                acl_data += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'

            acl_data += '</select>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('add_vote'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '' + \
                    '<form method="post">' + \
                        '<input name="name" placeholder="' + load_lang('name') + '">' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_100" name="subject" placeholder="' + load_lang('explanation') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<textarea class="opennamu_textarea_500" name="data" placeholder="' + load_lang('1_line_1_q') + '"></textarea>' + \
                        '<hr class="main_hr">' + \
                        '<input type="checkbox" value="Y" name="open_select"> ' + load_lang('open_vote') + \
                        '<h2>' + load_lang('period') + '</h2>'
                        '<input type="date" name="date" pattern="\\d{4}-\\d{2}-\\d{2}">' + \
                        '<hr class="main_hr">' + \
                        '<input type="checkbox" value="Y" name="limitless"> ' + load_lang('limitless') + \
                        '<h2>' + load_lang('acl') + '</h2>' + \
                        acl_data + ' <a href="/acl/TEST#exp">(' + load_lang('explanation') + ')</a>' + \
                        '<hr class="main_hr">' + \
                        '<button type="submit">' + load_lang('send') + '</buttom>' + \
                    '</form>' + \
                '',
                menu = [['vote', load_lang('return')]]
            ))