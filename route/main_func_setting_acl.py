from .tool.func import *

def main_func_setting_acl():
    with get_db_connect() as conn:
        curs = conn.cursor()

        i_list = {
            1 : 'edit',
            2 : 'discussion',
            3 : 'upload_acl',
            4 : 'all_view_acl',
            5 : 'many_upload_acl',
            6 : 'vote_acl'
        }

        if flask.request.method == 'POST':
            if admin_check(None, 'edit_set (acl)') != 1:
                return re_error('/ban')
            else:
                for i in i_list:
                    curs.execute(db_change("update other set data = ? where name = ?"), [
                        flask.request.form.get(i_list[i], 'normal'),
                        i_list[i]
                    ])

                conn.commit()

                return redirect('/setting/acl')
        else:
            d_list = {}

            if admin_check() != 1:
                disable = 'disabled'
            else:
                disable = ''

            for i in i_list:
                curs.execute(db_change('select data from other where name = ?'), [i_list[i]])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list[i] = sql_d[0][0]
                else:
                    curs.execute(db_change('insert into other (name, data) values (?, ?)'), [i_list[i], 'normal'])

                    d_list[i] = 'normal'

            conn.commit()

            acl_div = []
            for i in range(0, len(i_list)):
                acl_div += ['']

            acl_list = get_acl_list()
            for i in range(0, len(i_list)):
                for data_list in acl_list:
                    if data_list == d_list[i + 1]:
                        check = 'selected="selected"'
                    else:
                        check = ''

                    acl_div[i] += '<option value="' + data_list + '" ' + check + '>' + (data_list if data_list != '' else 'normal') + '</option>'

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('main_acl_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                        <hr class="main_hr">
                        <span>''' + load_lang('document_acl') + '''</span> 
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="edit">''' + acl_div[0] + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('discussion_acl') + '''</span>
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="discussion">''' + acl_div[1] + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('upload_acl') + '''</span>
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="upload_acl">''' + acl_div[2] + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('view_acl') + '''</span>
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="all_view_acl">''' + acl_div[3] + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('many_upload_acl') + '''</span>
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="many_upload_acl">''' + acl_div[4] + '''</select>
                        <hr class="main_hr">
                        <span>''' + load_lang('vote_acl') + '''</span>
                        <hr class="main_hr">
                        <select ''' + disable + ''' name="vote_acl">''' + acl_div[5] + '''</select>
                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))