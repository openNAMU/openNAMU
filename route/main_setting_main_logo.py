from .tool.func import *

async def main_setting_main_logo():
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)

        skin_list = [0] + load_skin(conn, '', 1)
        i_list = []
        for i in skin_list:
            i_list += [['logo', '' if i == 0 else i]]

        if flask.request.method == 'POST':
            for i in i_list:
                curs.execute(db_change("update other set data = ? where name = ? and coverage = ?"), [
                    flask.request.form.get(('main_css' if i[1] == '' else i[1]), ''),
                    i[0], 
                    i[1]
                ])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (logo)')

            return redirect(conn, '/setting/main/logo')
        else:
            d_list = []
            for i in i_list:
                curs.execute(db_change('select data from other where name = ? and coverage = ?'), [i[0], i[1]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list += [sql_d[0][0]]
                else:
                    curs.execute(db_change('insert into other (name, data, coverage) values (?, ?, ?)'), [i[0], '', i[1]])

                    d_list += ['']

            end_data = ''
            for i in range(0, len(skin_list)):
                end_data += '' + \
                    '<span>' + get_lang(conn, 'wiki_logo') + ' ' + ('(' + skin_list[i] + ')' if skin_list[i] != 0 else '') + ' (HTML)' + \
                    '<hr class="main_hr">' + \
                    '<input name="' + (skin_list[i] if skin_list[i] != 0 else 'main_css') + '" value="' + html.escape(d_list[i]) + '">' + \
                    '<hr class="main_hr">' + \
                ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'wiki_logo'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + end_data + '''
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['setting/main', get_lang(conn, 'return')]]
            ))