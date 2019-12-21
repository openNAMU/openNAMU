from .tool.func import *

def edit_many_delete_2(conn, app_var):
    curs = conn.cursor()

    ip = ip_check()
    if admin_check() != 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        all_title = re.findall(r'([^\n]+)\n', flask.request.form.get('content', '').replace('\r\n', '\n') + '\n')
        for name in all_title:
            curs.execute(db_change("select data from data where title = ?"), [name])
            data = curs.fetchall()
            if data:
                today = get_time()
                leng = '-' + str(len(data[0][0]))

                history_plus(
                    name,
                    '',
                    today,
                    ip,
                    flask.request.form.get('send', ''),
                    leng,
                    'delete'
                )

                curs.execute(db_change("select title, link from back where title = ? and not type = 'cat' and not type = 'no'"), [name])
                for data in curs.fetchall():
                    curs.execute(db_change("insert into back (title, link, type) values (?, ?, 'no')"), [data[0], data[1]])

                curs.execute(db_change("delete from back where link = ?"), [name])
                curs.execute(db_change("delete from data where title = ?"), [name])
                conn.commit()

            file_check = re.search('^file:(.+)\.(.+)$', name)
            if file_check:
                file_check = file_check.groups()
                os.remove(os.path.join(
                    app_var['path_data_image'],
                    hashlib.sha224(bytes(file_check[0], 'utf-8')).hexdigest() + '.' + file_check[1]
                ))

            curs.execute(db_change('select data from other where name = "count_all_title"'))
            curs.execute(db_change("update other set data = ? where name = 'count_all_title'"), [str(int(curs.fetchall()[0][0]) - 1)])

        return redirect('/recent_changes')
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('many_delete'), wiki_set(), custom(), other2([0, 0])],
            data = '''
                <form method="post">
                    <textarea rows="25" placeholder="''' + load_lang('many_delete_help') + '''" name="content"></textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('delete') + '''</button>
                </form>
            ''',
            menu = [['manager/1', load_lang('return')]]
        ))     