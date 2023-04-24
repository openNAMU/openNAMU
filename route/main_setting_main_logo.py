from .tool.func import *

def main_setting_main_logo():
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        if admin_check() != 1:
            return re_error('/ban')

        skin_list = [0] + load_skin('', 1)
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

            conn.commit()

            admin_check(None, 'edit_set (logo)')

            return redirect('/setting/main/logo')
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

            conn.commit()

            end_data = ''
            for i in range(0, len(skin_list)):
                end_data += '' + \
                    '<span>' + load_lang('wiki_logo') + ' ' + ('(' + skin_list[i] + ')' if skin_list[i] != 0 else '') + ' (HTML)' + \
                    '<hr class="main_hr">' + \
                    '<input name="' + (skin_list[i] if skin_list[i] != 0 else 'main_css') + '" value="' + html.escape(d_list[i]) + '">' + \
                    '<hr class="main_hr">' + \
                ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('wiki_logo'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + end_data + '''
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting/main', load_lang('return')]]
            ))