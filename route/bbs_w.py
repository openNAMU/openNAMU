from .tool.func import *

def bbs_w(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_id from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        if not curs.fetchall():
            return redirect('/bbs/main')

        bbs_num_str = str(bbs_num)

        data = ''
        data += '''
            <table id="main_table_set">
                <tr id="main_table_top_tr">
                    <td id="main_table_width">''' + load_lang('version') + '''</td>
                    <td id="main_table_width">''' + load_lang('editor') + '''</td>
                    <td id="main_table_width">''' + load_lang('time') + '''</td>
                </tr>
        '''
        
        temp_id = ''
        temp_dict = {}

        curs.execute(db_change('select set_name, set_data, set_code from bbs_data where set_id = ? order by set_code + 0 desc'), [bbs_num])
        db_data = curs.fetchall()
        for for_a in db_data + [['', '', '']]:
            if temp_id != for_a[2]:
                if temp_id != '':
                    data += '''
                        <tr>
                            <td>''' + temp_dict['code'] + '''</td>
                            <td>''' + ip_pas(temp_dict['user_id']) + '''</td>
                            <td>''' + temp_dict['date'] + '''</td>
                        </tr>
                        <tr>
                            <td colspan="3">''' + ('<a href="/bbs/w/' + bbs_num_str + '/' + temp_dict['code'] + '">' + temp_dict['title'] + '</a>') + '''</td>
                        </tr>
                    '''

                temp_id = for_a[2]
                temp_dict['code'] = for_a[2]

            temp_dict[for_a[0]] = for_a[1]

        data += '</table>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('bbs_main'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['bbs/main', load_lang('return')], ['bbs/edit/' + bbs_num_str, load_lang('add')], ['bbs/set/' + bbs_num_str, load_lang('setting')]]
        ))