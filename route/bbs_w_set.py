from .tool.func import *

def bbs_w_set(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_id from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        if not curs.fetchall():
            return redirect('/bbs/main')

        i_list = {
            1 : 'bbs_acl',
            2 : 'bbs_edit_acl',
            3 : 'bbs_comment_acl',
            4 : 'bbs_view_acl'
        }
        bbs_num_str = str(bbs_num)

        if flask.request.method == 'POST':
            if admin_check(None, 'bbs_set (acl)') != 1:
                return re_error('/ban')
            else:
                for i in i_list:
                    curs.execute(db_change("update bbs_set set set_data = ? where set_name = ? and set_id = ?"), [
                        flask.request.form.get(i_list[i], 'normal'),
                        i_list[i],
                        bbs_num
                    ])

                conn.commit()

                return redirect('/bbs/set/' + bbs_num_str)
        else:
            d_list = {}

            if admin_check() != 1:
                disable = 'disabled'
            else:
                disable = ''

            for i in i_list:
                curs.execute(db_change('select set_data from bbs_set where set_name = ? and set_id = ?'), [i_list[i], bbs_num])
                sql_d = curs.fetchall()
                if sql_d:
                    d_list[i] = sql_d[0][0]
                else:
                    curs.execute(db_change('insert into bbs_set (set_name, set_code, set_id, set_data) values (?, "", ?, ?)'), [i_list[i], bbs_num, 'normal'])
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
                data = render_simple_set('''
                    <form method="post">
                        <hr class="main_hr">
                        <a href="/acl/TEST#exp">(''' + load_lang('reference') + ''')</a>
                        
                        <h2>''' + load_lang('acl') + '''</h2>
                        <h3>''' + load_lang('bbs_view_acl') + '''</h3>
                        <select ''' + disable + ''' name="bbs_view_acl">''' + acl_div[3] + '''</select>

                        <h4>''' + load_lang('bbs_acl') + '''</h4>
                        <select ''' + disable + ''' name="bbs_acl">''' + acl_div[0] + '''</select>

                        <h5>''' + load_lang('bbs_edit_acl') + '''</h5>
                        <select ''' + disable + ''' name="bbs_edit_acl">''' + acl_div[1] + '''</select>

                        <h5>''' + load_lang('bbs_comment_acl') + '''</h5>
                        <select ''' + disable + ''' name="bbs_comment_acl">''' + acl_div[2] + '''</select>

                        <h2>''' + load_lang('markup') + '''</h2>
                        ''' + load_lang('not_working') + '''
                        
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str, load_lang('return')]]
            ))