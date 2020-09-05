from .tool.func import *

def vote_add_2(conn):
    curs = conn.cursor()

    if admin_check() != 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        vote_data = flask.request.form.get('data', 'test\ntest_2')
        if vote_data.count('\n') < 1:
            return re_error('/ban')

        curs.execute(db_change('select id from vote order by id + 0 desc limit 1'))
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
        conn.commit()

        return redirect('/vote')
    else:
        acl_data = '<select name="acl_select">'
        acl_list = get_acl_list()
        for data_list in acl_list:
            acl_data += '<option value="' + data_list + '">' + (data_list if data_list != '' else 'normal') + '</option>'

        acl_data += '</select>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('add_vote'), wiki_set(), custom(), other2([0, 0])],
            data = '' + \
                '<form method="post">' + \
                    '<input name="name" placeholder="' + load_lang('name') + '">' + \
                    '<hr class="main_hr">' + \
                    '<textarea rows="3" name="subject" placeholder="' + load_lang('explanation') + '"></textarea>' + \
                    '<hr class="main_hr">' + \
                    '<textarea rows="10" name="data" placeholder="' + load_lang('1_line_1_q') + '"></textarea>' + \
                    '<hr class="main_hr">' + \
                    '<input type="checkbox" value="Y" name="open_select"> ' + load_lang('open_vote') + \
                    '<h2>' + load_lang('acl') + '</h2>' + \
                    acl_data + ' <a href="/acl/TEST#exp">(' + load_lang('explanation') + ')</a>' + \
                    '<hr class="main_hr">' + \
                    '<button type="submit">' + load_lang('send') + '</buttom>' + \
                '</form>' + \
            '',
            menu = [['vote', load_lang('return')]]
        ))