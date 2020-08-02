from .tool.func import *

def vote_select_2(conn, num):
    curs = conn.cursor()

    curs.execute(db_change('select name, subject, data, type from vote where id = ? and user = ""'), [num])
    data_list = curs.fetchall()
    if not data_list:
        return redirect('/vote')

    if data_list[0][3] == 'close' or data_list[0][3] == 'n_close':
        return redirect('/end_vote/' + num)

    if acl_check('', 'vote', num) == 1:
        return redirect('/end_vote/' + num)

    curs.execute(db_change('select user from vote where id = ? and user = ?'), [num, ip_check()])
    if curs.fetchall():
        return redirect('/end_vote/' + num)

    vote_data = re.findall(r'([^\n]+)', data_list[0][2].replace('\r\n', '\n'))

    if flask.request.method == 'POST':
        try:
            vaild_check = int(flask.request.form.get('vote_data', '0'))
        except:
            return redirect('/vote/' + num)

        if len(vote_data) - 1 < vaild_check:
            return redirect('/vote/' + num)

        curs.execute(db_change("insert into vote (name, id, subject, data, user, type) values ('', ?, '', ?, ?, 'select')"), [
            num,
            str(vaild_check),
            ip_check()
        ])
        conn.commit()

        return redirect('/end_vote/' + num)
    else:
        data = '' + \
            '<h2>' + data_list[0][0] + '</h2>' + \
            '<b>' + data_list[0][1] + '</b>' + \
            '<hr class="main_hr">' + \
        ''

        select_data = '<select name="vote_data">'
        line_num = 0
        for i in vote_data:
            select_data += '<option value="' + str(line_num) + '">' + i + '</option>'
            line_num += 1

        select_data += '</select>'
        data += '' + \
            '<form method="post">' + \
                select_data + \
                '<hr class="main_hr">' + \
                '<button type="submit">' + load_lang('send') + '</buttom>' + \
            '</form>' + \
        ''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('vote'), wiki_set(), custom(), other2([0, 0])],
            data = data,
            menu = [['vote', load_lang('return')], ['end_vote/' + num, load_lang('result')]]
        ))